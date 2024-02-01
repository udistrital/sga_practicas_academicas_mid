package services

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_mid_practicas_academicas/models"
	"github.com/udistrital/utils_oas/request"
)

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
