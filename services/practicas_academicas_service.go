package services

import "fmt"

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
