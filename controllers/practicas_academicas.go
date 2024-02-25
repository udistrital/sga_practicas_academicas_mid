package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_practica_academica_mid/models"
	"github.com/udistrital/sga_practica_academica_mid/services"
	"github.com/udistrital/utils_oas/request"
)

// PracticasAcademicasController operations for Practicas_academicas
type PracticasAcademicasController struct {
	beego.Controller
}

// URLMapping ...
func (c *PracticasAcademicasController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("ConsultarInfoSolicitante", c.ConsultarInfoSolicitante)
	c.Mapping("ConsultarInfoColaborador", c.ConsultarInfoColaborador)
	c.Mapping("ConsultarParametros", c.ConsultarParametros)
	c.Mapping("ConsultarEspaciosAcademicos", c.ConsultarEspaciosAcademicos)
	c.Mapping("EnviarInvitaciones", c.EnviarInvitaciones)
}

// Post ...
// @Title Create
// @Description create Practicas_academicas
// @Param	body		body 	models.Practicas_academicas	true		"body for Practicas_academicas content"
// @Success 201 {object} models.Practicas_academicas
// @Failure 400 the request contains incorrect syntaxis
// @router / [post]
func (c *PracticasAcademicasController) Post() {
	var solicitud map[string]interface{}
	var resDocs []interface{}
	var Referencia string
	var SolicitudPost map[string]interface{}
	var SolicitantePost map[string]interface{}
	var SolicitudEvolucionEstadoPost map[string]interface{}
	var IdEstadoTipoSolicitud int
	resultado := make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &solicitud); err == nil {
		services.ManejoDocumentosPost(&solicitud, &resDocs)
		services.AsignaciónVariablesPost(&solicitud, &Referencia, resDocs, &IdEstadoTipoSolicitud)

		c.Data["json"] = services.SolicitudPracticasPost(IdEstadoTipoSolicitud, Referencia, solicitud, SolicitudPost, &resultado, &SolicitantePost, &SolicitudEvolucionEstadoPost, &alerta, &alertas, &errorGetAll)
	} else {
		services.ManejoError(&alerta, &alertas, "", &errorGetAll, err)
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	if !errorGetAll {
		services.ManejoExito(&alertas, &alerta, resultado)
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Practicas_academicas by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not found resource
// @router /:id [get]
func (c *PracticasAcademicasController) GetOne() {
	id_practica := c.Ctx.Input.Param(":id")
	var Solicitudes []map[string]interface{}
	var tipoSolicitud map[string]interface{}
	var Estados []map[string]interface{}
	var Comentario []map[string]interface{}
	resultado := make(map[string]interface{})
	var errorGetAll bool

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante?query=SolicitudId.Id:"+id_practica, &Solicitudes)
	if errSolicitud == nil {
		if Solicitudes != nil && fmt.Sprintf("%v", Solicitudes[0]) != "map[]" {
			services.ManejoSolicitudesGetOne(Solicitudes, id_practica, &resultado, &tipoSolicitud, &Estados, &Comentario)
		} else {
			errorGetAll = true
			c.Data["message"] = "Error service GetAll: No data found"
			c.Abort("404")
		}
	} else {
		errorGetAll = true
		c.Data["message"] = "Error service GetAll: " + errSolicitud.Error()
		c.Abort("400")
	}

	if !errorGetAll {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}

	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Practicas_academicas
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not data found
// @Failure 400 the request contains incorrect syntax
// @router / [get]
func (c *PracticasAcademicasController) GetAll() {
	var query string
	var fields string
	var Solicitudes []map[string]interface{}
	var TipoEstado map[string]interface{}
	resultado := []interface{}{}
	var errorGetAll bool

	// query: k:v,k:v
	if query = c.GetString("query"); query != "" {
		query = "&query=" + query
	}
	// fields: col1,col2,entity.col3
	if fields = c.GetString("fields"); fields != "" {
		fields = "&fields=" + fields
	}

	c.Data["json"] = services.SolicitudGetAllSolicitudes(query, Solicitudes, &TipoEstado, &resultado, errorGetAll)

	if !errorGetAll {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
	}

	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Practicas_academicas
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Practicas_academicas	true		"body for Practicas_academicas content"
// @Success 200 {object} models.Practicas_academicas
// @Failure 403 :id is not int
// @router /:id [put]
func (c *PracticasAcademicasController) Put() {
	id_practica := c.Ctx.Input.Param(":id")
	var RespuestaSolicitud map[string]interface{}
	var Solicitud map[string]interface{}
	var SolicitudPut map[string]interface{}
	var NuevoEstado map[string]interface{}
	var anteriorEstado []map[string]interface{}
	var tipoSolicitud map[string]interface{}

	var SolicitudEvolucionEstadoPost map[string]interface{}
	var anteriorEstadoPost map[string]interface{}
	var observacionPost map[string]interface{}
	var Referencia string
	var resDocs []interface{}
	var resultado = make(map[string]interface{})
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &RespuestaSolicitud); err == nil {

		// Consulta de información de la solicitud
		dataMessage, dataJson, exito := services.ManejoSolicitudesPut(&NuevoEstado, id_practica, RespuestaSolicitud, &SolicitudEvolucionEstadoPost, &SolicitudPut, &Solicitud, &resultado, &errorGetAll, &alertas, &alerta, &Referencia, &observacionPost, &tipoSolicitud, &anteriorEstado, &anteriorEstadoPost, &resDocs)

		if !exito {
			if dataMessage != nil {
				c.Data["message"] = dataMessage
				c.Abort("404")
			} else {
				c.Data["json"] = dataJson
			}
		}
	} else {
		services.ManejoError(&alerta, &alertas, "", &errorGetAll, err)
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	if !errorGetAll {
		services.ManejoExito(&alertas, &alerta, resultado)
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}

// ConsultarInfoSolicitante ...
// @Title ConsultarInfoSolicitante
// @Description get información del docente solicitante de la practica academica
// @Param	id		id perteneciente a terceros
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not found resource
// @router /consultar_solicitante/:id [get]
func (c *PracticasAcademicasController) ConsultarInfoSolicitante() {
	idTercero := c.Ctx.Input.Param(":id")

	var resultado = make(map[string]interface{})
	var persona []map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{}
	var errorGetAll bool

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero?query=Id:"+idTercero, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		var tipoVinculacion []map[string]interface{}
		var correoElectronico []map[string]interface{}
		var correoInstitucional []map[string]interface{}
		var correoPersonal []map[string]interface{}
		var telefono []map[string]interface{}
		var celular []map[string]interface{}
		var jsondata map[string]interface{}

		// Correo institucional --> 94
		c.Data["json"] = services.SolicitudCorreoInstitucionalConsultarInfo(idTercero, &correoInstitucional, &jsondata, &resultado, &alertas, &alerta, &errorGetAll)

		// Correo --> 53
		c.Data["json"] = services.SolicitudCorreoConsultarInfo(idTercero, &correoElectronico, &jsondata, &resultado, &alertas, &alerta, &errorGetAll)

		// Correo personal --> 253
		c.Data["json"] = services.SolicitudCorreoPersonalConsultarInfo(idTercero, &correoPersonal, &jsondata, &resultado, &alertas, &alerta, &errorGetAll)

		// Teléfono --> 51
		c.Data["json"] = services.SolicitudTelefonoConsultarInfo(idTercero, &telefono, &jsondata, &resultado, &alertas, &alerta, &errorGetAll)

		// Celular --> 52
		c.Data["json"] = services.SolicitudCelularConsultarInfo(idTercero, &celular, &jsondata, &resultado, &alertas, &alerta, &errorGetAll)

		// DOCENTE DE PLANTA 	292
		// DOCENTE DE CARRERA TIEMPO COMPLETO 	293
		// DOCENTE DE CARRERA MEDIO TIEMPO 	294
		// DOCENTE DE VINCULACIÓN ESPECIAL 	295
		// HORA CÁTEDRA 	297
		// TIEMPO COMPLETO OCASIONAL 	296
		// MEDIO TIEMPO OCASIONAL 	298
		// HORA CÁTEDRA POR HONORARIOS 	299
		c.Data["json"] = services.SolicitudTipoVinculacionConsultarInfo(idTercero, &tipoVinculacion, &resultado, &alertas, &alerta, &errorGetAll)

		resultado["Nombre"] = persona[0]["NombreCompleto"]
		resultado["Id"], _ = strconv.ParseInt(idTercero, 10, 64)
		resultado["PuedeBorrar"] = false

		c.Data["json"] = resultado
	} else {
		logs.Error(errPersona)
		errorGetAll = true
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Data not found", "Data": nil}
	}

	if !errorGetAll {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
	}

	c.ServeJSON()
}

// ConsultarInfoColaborador ...
// @Title ConsultarInfoColaborador
// @Description get información del docente colaborador
// @Param	id		documento de identidad del usuario registrado en wso2
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not found resource
// @router /consultar_colaborador/:id [get]
func (c *PracticasAcademicasController) ConsultarInfoColaborador() {
	idStr := c.Ctx.Input.Param(":id")
	var resultado = make(map[string]interface{})
	var persona []map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{}
	var errorGetAll bool

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Numero:"+idStr, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		if persona[0]["Status"] != 404 {
			var tipoVinculacion []map[string]interface{}
			var correoElectronico []map[string]interface{}
			var correoInstitucional []map[string]interface{}
			var correoPersonal []map[string]interface{}
			var telefono []map[string]interface{}
			var celular []map[string]interface{}
			var jsondata map[string]interface{}

			idTercero := persona[0]["TerceroId"].(map[string]interface{})["Id"]

			// DOCENTE DE PLANTA 	292
			// DOCENTE DE CARRERA TIEMPO COMPLETO 	293
			// DOCENTE DE CARRERA MEDIO TIEMPO 	294
			// DOCENTE DE VINCULACIÓN ESPECIAL 	295
			// HORA CÁTEDRA 	297
			// TIEMPO COMPLETO OCASIONAL 	296
			// MEDIO TIEMPO OCASIONAL 	298
			// HORA CÁTEDRA POR HONORARIOS 	299
			c.Data["json"] = services.SolicitudTipoVinculacionInfoColaborador(&tipoVinculacion, idTercero, &jsondata, &resultado, &alertas, &alerta, &errorGetAll, &correoInstitucional, &correoElectronico, &correoPersonal, &telefono, &celular, persona)
		} else {
			if persona[0]["Message"] == "Not found resource" {
				c.Data["json"] = nil
			} else {
				services.ManejoError(&alerta, &alertas, "No data found", &errorGetAll)
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Data not found", "Data": nil}
			}
		}
	} else {
		logs.Error(errPersona)
		errorGetAll = true
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Data not found", "Data": nil}
	}

	if !errorGetAll {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": resultado}
	}

	c.ServeJSON()
}

// ConsultarParametros ...
// @Title ConsultarParametros
// @Description get parametros para creación de practica academica
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not found resource
// @router /consultar_parametros/ [get]
func (c *PracticasAcademicasController) ConsultarParametros() {
	var getProyecto []map[string]interface{}
	var proyectos []map[string]interface{}
	var estados []interface{}
	var vehiculos map[string]interface{}
	var resultado = make(map[string]interface{})
	var periodos map[string]interface{}
	var tipoEstados map[string]interface{}
	var tipoSolicitud map[string]interface{}

	if errPeriodo := services.SolicitudPeriodoParametros(&periodos, &resultado); errPeriodo != nil {
		c.Data["system"] = errPeriodo
		c.Abort("404")
	}

	if errProyecto := services.SolicitudProyectoParametros(&getProyecto, &resultado, proyectos); errProyecto != nil {
		c.Data["system"] = errProyecto
		c.Abort("404")
	}

	if errVehiculo := services.SolicitudVehiculoParametros(&vehiculos, &resultado, getProyecto); errVehiculo != nil {
		c.Data["system"] = errVehiculo
		c.Abort("404")
	}

	if errTipoSolicitud := services.SolicitudTipoParametros(&tipoSolicitud, &tipoEstados, &resultado, estados); errTipoSolicitud != nil {
		c.Data["system"] = errTipoSolicitud
		c.Abort("404")
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": resultado}

	c.ServeJSON()
}

// ConsultarEspaciosAcademicos ...
// @Title ConsultarEspaciosAcademicos
// @Description get estados de practica academica
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not found resource
// @router /consultar_espacios_academicos/:id [get]
func (c *PracticasAcademicasController) ConsultarEspaciosAcademicos() {
	resultado := []interface{}{}
	var espaciosAcademicos map[string]interface{}
	idStr := c.Ctx.Input.Param(":id")

	errEspaciosAcademicos := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico?query=activo:true,docente_id:"+fmt.Sprintf("%v", idStr), &espaciosAcademicos)
	if errEspaciosAcademicos == nil && fmt.Sprintf("%v", espaciosAcademicos["Data"]) != "[map[]]" {
		services.AsignarResultadoEspaciosAcademicos(&resultado, espaciosAcademicos)
	} else {
		resultado = nil
		logs.Error(espaciosAcademicos)
		c.Data["system"] = errEspaciosAcademicos
		c.Abort("404")
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": resultado}

	c.ServeJSON()
}

// EnviarInvitaciones ...
// @Title EnviarInvitaciones
// @Description enviar invitaciones al correo de los estudiantes
// @Param	body		body 	models.Practicas_academicas	true		"body for Practicas_academicas content"
// @Success 201 {object} models.Practicas_academicas
// @Failure 400 the request contains incorrect syntaxis
// @router /enviar_invitacion/ [post]
func (c *PracticasAcademicasController) EnviarInvitaciones() {

	var Solicitudes []map[string]interface{}
	var CorreoPost map[string]interface{}
	var solicitud map[string]interface{}
	var alerta models.Alert
	var errorGetAll bool
	alertas := []interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &solicitud); err == nil {
		id_practica := solicitud["Id"]

		//*********************//
		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante?query=SolicitudId.Id:"+id_practica.(string), &Solicitudes)
		if errSolicitud == nil {
			if Solicitudes != nil && fmt.Sprintf("%v", Solicitudes[0]) != "map[]" {
				Referencia := Solicitudes[0]["SolicitudId"].(map[string]interface{})["Referencia"].(string)
				var ReferenciaJson map[string]interface{}
				if err := json.Unmarshal([]byte(Referencia), &ReferenciaJson); err == nil {
					ReferenciaJson["Id"] = id_practica
				}

				idEstado := fmt.Sprintf("%v", Solicitudes[0]["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["Id"].(float64))

				c.Data["json"] = services.ValidarEstadoEnviarInvitaciones(idEstado, ReferenciaJson, &CorreoPost, &errorGetAll, &alertas, &alerta)

			} else {
				errorGetAll = true
				c.Data["message"] = "Error service GetAll: No data found"
				c.Abort("404")
			}
		} else {
			errorGetAll = true
			c.Data["message"] = "Error service GetAll: " + errSolicitud.Error()
			c.Abort("400")
		}
		//*********************//

	}

	if !errorGetAll {
		alertas = append(alertas, "Correos enviados")
		alerta.Code = "200"
		alerta.Type = "OK"
		alerta.Body = alertas
		c.Data["json"] = map[string]interface{}{"Response": alerta}
	}

	c.ServeJSON()
}
