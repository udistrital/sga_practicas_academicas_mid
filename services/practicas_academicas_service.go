package services

import "fmt"

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
