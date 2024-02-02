package services

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid_practicas_academicas/models"
	"github.com/udistrital/utils_oas/request"
)

// FUNCIONES QUE SE USAN EN GETONE

func SolicitudTipoGetOne(Solicitudes []map[string]interface{}, tipoSolicitud map[string]interface{}, resultado *map[string]interface{}) {
	idEstado := fmt.Sprintf("%v", Solicitudes[0]["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["Id"].(float64))

	errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado_tipo_solicitud?query=Activo:true,Id:"+idEstado, &tipoSolicitud)
	if errTipoSolicitud == nil {
		if tipoSolicitud != nil && fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0]) != "map[]" {
			(*resultado)["EstadoTipoSolicitudId"] = tipoSolicitud["Data"].([]interface{})[0]
		}
	}
}

func ManejoEstadosGetOne(id_practica string, Estados []map[string]interface{}, Comentario []map[string]interface{}, resultado *map[string]interface{}) {
	errEstados := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado?query=SolicitudId.Id:"+id_practica, &Estados)
	if errEstados == nil {
		if Estados != nil && fmt.Sprintf("%v", Estados[0]) != "map[]" {
			for _, v := range Estados {

				errComentario := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"observacion?query=titulo:"+fmt.Sprintf("%v", v["Id"]), &Comentario)
				if errComentario == nil {
					if Comentario != nil && fmt.Sprintf("%v", Comentario[0]) != "map[]" {
						v["Comentario"] = Comentario[0]["Valor"]
					} else {
						v["Comentario"] = ""
					}
				}
			}
			(*resultado)["Estados"] = Estados
		}
	}
}

func ManejoSolicitudesGetOne(Solicitudes []map[string]interface{}, id_practica string, resultado *map[string]interface{}, tipoSolicitud map[string]interface{}, Estados []map[string]interface{}, Comentario []map[string]interface{}) {
	Referencia := Solicitudes[0]["SolicitudId"].(map[string]interface{})["Referencia"].(string)
	fechaRadicado := Solicitudes[0]["SolicitudId"].(map[string]interface{})["FechaRadicacion"].(string)
	var ReferenciaJson map[string]interface{}
	if err := json.Unmarshal([]byte(Referencia), &ReferenciaJson); err == nil {
		ReferenciaJson["Id"] = id_practica
		*resultado = ReferenciaJson
		(*resultado)["FechaRadicado"] = fechaRadicado
	}

	SolicitudTipoGetOne(Solicitudes, tipoSolicitud, resultado)
	ManejoEstadosGetOne(id_practica, Estados, Comentario, resultado)
}

// FUNCIONES QUE SE USAN EN CONSULTAR INFO SOLICITANTE

func ManejoVinculacionSolicitante(resultado *map[string]interface{}, tipoVinculacion []map[string]interface{}) {
	for _, tv := range tipoVinculacion {
		if fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "292" ||
			fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "293" ||
			fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "294" ||
			fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "295" ||
			fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "296" ||
			fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "297" ||
			fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "298" ||
			fmt.Sprintf("%v", tv["TipoVinculacionId"]) == "299" {

			var vinculacion map[string]interface{}
			errVinculacion := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro?query=Id:"+fmt.Sprintf("%v", tv["TipoVinculacionId"])+",Activo:true&limit=0", &vinculacion)
			if errVinculacion == nil && fmt.Sprintf("%v", vinculacion["Data"]) != "[map[]]" {
				if vinculacion["Status"] != 404 {
					(*resultado)["TipoVinculacionId"] = vinculacion["Data"].([]interface{})[0]
				}
			}
		}
	}
}

// FUNCIONES QUE SE USAN EN CONSULTAR PARAMETROS

func ManejoProyectosParametros(resultado *map[string]interface{}, getProyecto *[]map[string]interface{}, proyectos []map[string]interface{}) {
	for _, proyectoAux := range *getProyecto {
		proyecto := map[string]interface{}{
			"Id":          proyectoAux["Id"],
			"Nombre":      proyectoAux["Nombre"],
			"Codigo":      proyectoAux["Codigo"],
			"CodigoSnies": proyectoAux["CodigoSnies"],
		}
		proyectos = append(proyectos, proyecto)
	}
	(*resultado)["proyectos"] = proyectos
}

func ManejoEstadosParametros(resultado *map[string]interface{}, tipoEstados map[string]interface{}, estados []interface{}) {
	if tipoEstados["Status"] != "404" {
		for _, estado := range tipoEstados["Data"].([]interface{}) {
			estados = append(estados, estado.(map[string]interface{})["EstadoId"])
			(*resultado)["estados"] = estados
		}
	}
}

// FUNCIONES QUE SE USAN EN CONSULTAR ESPACIOS ACADEMICOS

func AsignarResultadoEspaciosAcademicos(resultado *[]interface{}, espaciosAcademicos map[string]interface{}) {
	if espaciosAcademicos["Status"] != "404" {

		for _, espacioAcademico := range espaciosAcademicos["Data"].([]interface{}) {

			*resultado = append(*resultado, map[string]interface{}{
				"Nombre": fmt.Sprintf("%v", espacioAcademico.(map[string]interface{})["nombre"]) + " - " + fmt.Sprintf("%v", espacioAcademico.(map[string]interface{})["grupo"]),
				"Id":     espacioAcademico.(map[string]interface{})["_id"],
			})
		}
	}
}

// FUNCIONES QUE SE USAN EN VARIOS ENDPOINTS

func ManejoError(alerta *models.Alert, alertas *[]interface{}, mensaje string, err ...error) {
	var msj string
	if len(err) > 0 && err[0] != nil {
		msj = mensaje + err[0].Error()
	} else {
		msj = mensaje
	}
	*alertas = append(*alertas, msj)
	(*alerta).Body = *alertas
	(*alerta).Type = "error"
	(*alerta).Code = "400"
}
