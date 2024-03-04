package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid_practicas_academicas/services"
	"github.com/udistrital/utils_oas/errorhandler"
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
	defer errorhandler.HandlePanic(&c.Controller)
	dataBody := c.Ctx.Input.RequestBody
	resultado := services.CrearPracticaAcademica(dataBody)
	c.Data["json"] = resultado
	c.Ctx.Output.SetStatus(resultado.Status)
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
	defer errorhandler.HandlePanic(&c.Controller)
	id_practica := c.Ctx.Input.Param(":id")
	resultado := services.GetPracticaAcademica(id_practica)
	c.Data["json"] = resultado
	c.Ctx.Output.SetStatus(resultado.Status)
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
	defer errorhandler.HandlePanic(&c.Controller)
	var query string = ""
	var fields string = ""
	// query: k:v,k:v
	if query = c.GetString("query"); query != "" {
		query = "&query=" + query
	}
	// fields: col1,col2,entity.col3
	if fields = c.GetString("fields"); fields != "" {
		fields = "&fields=" + fields
	}
	resultado := services.GetAllPracticasAcademicas(query, fields)
	c.Data["json"] = resultado
	c.Ctx.Output.SetStatus(resultado.Status)
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
	defer errorhandler.HandlePanic(&c.Controller)
	id_practica := c.Ctx.Input.Param(":id")
	dataBody := c.Ctx.Input.RequestBody
	resultado := services.ActualizarPracticaAcademica(id_practica, dataBody)
	c.Data["json"] = resultado
	c.Ctx.Output.SetStatus(resultado.Status)
	c.ServeJSON()
}

// ConsultarInfoSolicitante ...
// @Title ConsultarInfoSolicitante
// @Description get información del docente solicitante de la practica academica
// @Param	id		id perteneciente a terceros
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not found resource
// @router /solicitantes/:id [get]
func (c *PracticasAcademicasController) ConsultarInfoSolicitante() {
	defer errorhandler.HandlePanic(&c.Controller)
	idTercero := c.Ctx.Input.Param(":id")
	resultado := services.ConsultarInfoSolicitante(idTercero)
	c.Data["json"] = resultado
	c.Ctx.Output.SetStatus(resultado.Status)
	c.ServeJSON()
}

// ConsultarInfoColaborador ...
// @Title ConsultarInfoColaborador
// @Description get información del docente colaborador
// @Param	id		documento de identidad del usuario registrado en wso2
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not found resource
// @router /colaboradores/:id [get]
func (c *PracticasAcademicasController) ConsultarInfoColaborador() {
	defer errorhandler.HandlePanic(&c.Controller)
	idStr := c.Ctx.Input.Param(":id")
	resultado := services.ConsultarInfoColaborador(idStr)
	c.Data["json"] = resultado
	c.Ctx.Output.SetStatus(resultado.Status)
	c.ServeJSON()
}

// ConsultarParametros ...
// @Title ConsultarParametros
// @Description get parametros para creación de practica academica
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not found resource
// @router /parametros/ [get]
func (c *PracticasAcademicasController) ConsultarParametros() {
	defer errorhandler.HandlePanic(&c.Controller)
	resultado := services.ConsultarParametros()
	c.Data["json"] = resultado
	c.Ctx.Output.SetStatus(resultado.Status)
	c.ServeJSON()
}

// ConsultarEspaciosAcademicos ...
// @Title ConsultarEspaciosAcademicos
// @Description get estados de practica academica
// @Success 200 {object} models.Practicas_academicas
// @Failure 404 not found resource
// @router /espacios-academicos/:id [get]
func (c *PracticasAcademicasController) ConsultarEspaciosAcademicos() {
	defer errorhandler.HandlePanic(&c.Controller)
	idStr := c.Ctx.Input.Param(":id")
	resultado := services.ConsultarEspaciosAcademicos(idStr)
	c.Data["json"] = resultado
	c.Ctx.Output.SetStatus(resultado.Status)
	c.ServeJSON()
}

// EnviarInvitaciones ...
// @Title EnviarInvitaciones
// @Description enviar invitaciones al correo de los estudiantes
// @Param	body		body 	models.Practicas_academicas	true		"body for Practicas_academicas content"
// @Success 201 {object} models.Practicas_academicas
// @Failure 400 the request contains incorrect syntaxis
// @router /enviar-invitacion/ [post]
func (c *PracticasAcademicasController) EnviarInvitaciones() {
	defer errorhandler.HandlePanic(&c.Controller)
	dataBody := c.Ctx.Input.RequestBody
	resultado := services.EnviarInvitaciones(dataBody)
	c.Data["json"] = resultado
	c.Ctx.Output.SetStatus(resultado.Status)
	c.ServeJSON()
}
