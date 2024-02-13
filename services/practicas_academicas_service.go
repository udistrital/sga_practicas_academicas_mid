package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid_practicas_academicas/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

// FUNCIONES QUE SE USAN EN POST

func ManejoDocumentosPost(solicitud *map[string]interface{}, resDocs *[]interface{}) {
	for i := range (*solicitud)["Documentos"].([]interface{}) {
		auxDoc := []map[string]interface{}{}
		documento := map[string]interface{}{
			"IdTipoDocumento": (*solicitud)["Documentos"].([]interface{})[i].(map[string]interface{})["IdTipoDocumento"],
			"nombre":          (*solicitud)["Documentos"].([]interface{})[i].(map[string]interface{})["nombre"],
			"metadatos":       (*solicitud)["Documentos"].([]interface{})[i].(map[string]interface{})["metadatos"],
			"descripcion":     (*solicitud)["Documentos"].([]interface{})[i].(map[string]interface{})["descripcion"],
			"file":            (*solicitud)["Documentos"].([]interface{})[i].(map[string]interface{})["file"],
		}
		auxDoc = append(auxDoc, documento)
		doc, errDoc := models.RegistrarDoc(auxDoc)
		if errDoc == nil {
			docTem := map[string]interface{}{
				"Nombre":        doc.(map[string]interface{})["Nombre"].(string),
				"Enlace":        doc.(map[string]interface{})["Enlace"],
				"Id":            doc.(map[string]interface{})["Id"],
				"TipoDocumento": doc.(map[string]interface{})["TipoDocumento"],
				"Activo":        doc.(map[string]interface{})["Activo"],
			}

			*resDocs = append(*resDocs, docTem)
		}
	}
}

func AsignaciónVariablesPost(solicitud *map[string]interface{}, Referencia *string, resDocs []interface{}, IdEstadoTipoSolicitud *int) {
	jsonPeriodo, _ := json.Marshal((*solicitud)["Periodo"])
	jsonDocumento, _ := json.Marshal(resDocs)
	jsonProyecto, _ := json.Marshal((*solicitud)["Proyecto"])
	jsonEspacio, _ := json.Marshal((*solicitud)["EspacioAcademico"])
	jsonVehiculo, _ := json.Marshal((*solicitud)["TipoVehiculo"])
	jsonDocente, _ := json.Marshal((*solicitud)["DocenteSolicitante"])
	jsonDocentes, _ := json.Marshal((*solicitud)["DocentesInvitados"])

	*Referencia = "{\n\"Periodo\":" + fmt.Sprintf("%v", string(jsonPeriodo)) +
		",\n\"Proyecto\": " + fmt.Sprintf("%v", string(jsonProyecto)) +
		",\n\"EspacioAcademico\": " + fmt.Sprintf("%v", string(jsonEspacio)) +
		",\n\"Semestre\": " + fmt.Sprintf("%v", (*solicitud)["Semestre"]) +
		",\n\"NumeroEstudiantes\": " + fmt.Sprintf("%v", (*solicitud)["NumeroEstudiantes"]) +
		",\n\"NumeroGrupos\": " + fmt.Sprintf("%v", (*solicitud)["NumeroGrupos"]) +
		",\n\"Duracion\": " + fmt.Sprintf("%v", (*solicitud)["Duracion"]) +
		",\n\"NumeroVehiculos\": " + fmt.Sprintf("%v", (*solicitud)["NumeroVehiculos"]) +
		",\n\"TipoVehiculo\": " + fmt.Sprintf("%v", string(jsonVehiculo)) +
		",\n\"FechaHoraSalida\": \"" + time_bogota.TiempoCorreccionFormato((*solicitud)["FechaHoraSalida"].(string)) + "\"" +
		",\n\"FechaHoraRegreso\": \"" + time_bogota.TiempoCorreccionFormato((*solicitud)["FechaHoraRegreso"].(string)) + "\"" +
		",\n\"Documentos\": " + fmt.Sprintf("%v", string(jsonDocumento)) +
		",\n\"DocenteSolicitante\": " + fmt.Sprintf("%v", string(jsonDocente)) +
		",\n\"DocentesInvitados\": " + fmt.Sprintf("%v", string(jsonDocentes)) + "\n}"

	*IdEstadoTipoSolicitud = 34
}

func solicitudEliminarRegistroPost(IdSolicitud interface{}, alerta *models.Alert, alertas *[]interface{}, errorGetAll *bool, errSolicitante error) interface{} {
	var resultado2 map[string]interface{}
	request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+fmt.Sprintf("%v", IdSolicitud), "DELETE", &resultado2, nil)
	ManejoError(alerta, alertas, "", errorGetAll, errSolicitante)
	return map[string]interface{}{"Response": *alerta}
}

func solicitudEliminarSolicitantePost(IdSolicitud interface{}, SolicitantePost map[string]interface{}, alerta *models.Alert, alertas *[]interface{}, errorGetAll *bool, errSolicitante error) interface{} {
	var resultado2 map[string]interface{}
	request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+fmt.Sprintf("%v", IdSolicitud), "DELETE", &resultado2, nil)
	request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante/"+fmt.Sprintf("%v", SolicitantePost["Id"]), "DELETE", &resultado2, nil)
	ManejoError(alerta, alertas, "", errorGetAll, errSolicitante)
	return map[string]interface{}{"Response": *alerta}
}

func solicitudEvolucionEstadoPost(solicitud map[string]interface{}, IdSolicitud interface{}, IdEstadoTipoSolicitud int, SolicitudEvolucionEstadoPost *map[string]interface{}, resultado *map[string]interface{}, SolicitantePost map[string]interface{}, alerta *models.Alert, alertas *[]interface{}, errorGetAll *bool, errSolicitante error) interface{} {
	SolicitudEvolucionEstado := map[string]interface{}{
		"TerceroId": solicitud["SolicitanteId"],
		"SolicitudId": map[string]interface{}{
			"Id": IdSolicitud,
		},
		"EstadoTipoSolicitudIdAnterior": nil,
		"EstadoTipoSolicitudId": map[string]interface{}{
			"Id": IdEstadoTipoSolicitud,
		},
		"Activo":      true,
		"FechaLimite": fmt.Sprintf("%v", solicitud["FechaRadicacion"]),
	}

	errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado", "POST", SolicitudEvolucionEstadoPost, SolicitudEvolucionEstado)
	if errSolicitudEvolucionEstado == nil {
		if *SolicitudEvolucionEstadoPost != nil && fmt.Sprintf("%v", *SolicitudEvolucionEstadoPost) != "map[]" {
			(*resultado)["Solicitante"] = SolicitantePost["Data"]
			return nil
		} else {
			ManejoError(alerta, alertas, "No data found", errorGetAll)
			return map[string]interface{}{"Response": *alerta}
		}
	} else {
		return solicitudEliminarSolicitantePost(IdSolicitud, SolicitantePost, alerta, alertas, errorGetAll, errSolicitante)
	}
}

func solicitudTablaSolicitantePost(solicitud map[string]interface{}, IdSolicitud interface{}, SolicitantePost *map[string]interface{}, IdEstadoTipoSolicitud int, SolicitudEvolucionEstadoPost *map[string]interface{}, resultado *map[string]interface{}, alerta *models.Alert, alertas *[]interface{}, errorGetAll *bool) interface{} {
	Solicitante := map[string]interface{}{
		"TerceroId": solicitud["SolicitanteId"],
		"SolicitudId": map[string]interface{}{
			"Id": IdSolicitud,
		},
		"Activo": true,
	}

	errSolicitante := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante", "POST", SolicitantePost, Solicitante)
	if errSolicitante == nil && fmt.Sprintf("%v", (*SolicitantePost)["Status"]) != "400" {
		if *SolicitantePost != nil && fmt.Sprintf("%v", *SolicitantePost) != "map[]" {
			//POST a la tabla solicitud_evolucion estado
			return solicitudEvolucionEstadoPost(solicitud, IdSolicitud, IdEstadoTipoSolicitud, SolicitudEvolucionEstadoPost, resultado, *SolicitantePost, alerta, alertas, errorGetAll, errSolicitante)
		} else {
			ManejoError(alerta, alertas, "No data found", errorGetAll)
			return map[string]interface{}{"Response": *alerta}
		}
	} else {
		//Se elimina el registro de solicitud si no se puede hacer el POST a la tabla solicitante
		return solicitudEliminarRegistroPost(IdSolicitud, alerta, alertas, errorGetAll, errSolicitante)
	}
}

func SolicitudPracticasPost(IdEstadoTipoSolicitud int, Referencia string, solicitud map[string]interface{}, SolicitudPost map[string]interface{}, resultado *map[string]interface{}, SolicitantePost *map[string]interface{}, SolicitudEvolucionEstadoPost *map[string]interface{}, alerta *models.Alert, alertas *[]interface{}, errorGetAll *bool) interface{} {
	SolicitudPracticas := map[string]interface{}{
		"EstadoTipoSolicitudId": map[string]interface{}{"Id": IdEstadoTipoSolicitud},
		"Referencia":            Referencia,
		"Resultado":             "",
		"FechaRadicacion":       fmt.Sprintf("%v", solicitud["FechaRadicacion"]),
		"Activo":                true,
		"SolicitudPadreId":      nil,
	}

	errSolicitud := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud", "POST", &SolicitudPost, SolicitudPracticas)
	if errSolicitud == nil {
		if SolicitudPost["Success"] != false && fmt.Sprintf("%v", SolicitudPost) != "map[]" {
			(*resultado)["Solicitud"] = SolicitudPost["Data"]
			IdSolicitud := SolicitudPost["Data"].(map[string]interface{})["Id"]

			//POST tabla solicitante
			return solicitudTablaSolicitantePost(solicitud, IdSolicitud, SolicitantePost, IdEstadoTipoSolicitud, SolicitudEvolucionEstadoPost, resultado, alerta, alertas, errorGetAll)
		} else {
			ManejoError(alerta, alertas, "No data found"+fmt.Sprintf("%v", SolicitudPracticas), errorGetAll)
			return map[string]interface{}{"Response": *alerta}
		}
	} else {
		ManejoError(alerta, alertas, "", errorGetAll, errSolicitud)
		return map[string]interface{}{"Response": *alerta}
	}
}

// FUNCIONES QUE SE USAN EN GETONE

func SolicitudTipoGetOne(Solicitudes []map[string]interface{}, tipoSolicitud *map[string]interface{}, resultado *map[string]interface{}) {
	idEstado := fmt.Sprintf("%v", Solicitudes[0]["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["Id"].(float64))

	errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado_tipo_solicitud?query=Activo:true,Id:"+idEstado, tipoSolicitud)
	if errTipoSolicitud == nil {
		if *tipoSolicitud != nil && fmt.Sprintf("%v", (*tipoSolicitud)["Data"].([]interface{})[0]) != "map[]" {
			(*resultado)["EstadoTipoSolicitudId"] = (*tipoSolicitud)["Data"].([]interface{})[0]
		}
	}
}

func ManejoEstadosGetOne(id_practica string, Estados *[]map[string]interface{}, Comentario *[]map[string]interface{}, resultado *map[string]interface{}) {
	errEstados := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado?query=SolicitudId.Id:"+id_practica, Estados)
	if errEstados == nil {
		if *Estados != nil && fmt.Sprintf("%v", (*Estados)[0]) != "map[]" {
			for _, v := range *Estados {

				errComentario := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"observacion?query=titulo:"+fmt.Sprintf("%v", v["Id"]), Comentario)
				if errComentario == nil {
					if *Comentario != nil && fmt.Sprintf("%v", (*Comentario)[0]) != "map[]" {
						v["Comentario"] = (*Comentario)[0]["Valor"]
					} else {
						v["Comentario"] = ""
					}
				}
			}
			(*resultado)["Estados"] = *Estados
		}
	}
}

func ManejoSolicitudesGetOne(Solicitudes []map[string]interface{}, id_practica string, resultado *map[string]interface{}, tipoSolicitud *map[string]interface{}, Estados *[]map[string]interface{}, Comentario *[]map[string]interface{}) {
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

// FUNCIONES QUE SE USAN EN GETALL

func manejoSolicitudesGetAll(Solicitudes []map[string]interface{}, TipoEstado *map[string]interface{}, resultado *[]interface{}) {
	for _, solicitud := range Solicitudes {
		errTipoEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado_tipo_solicitud?query=Id:"+fmt.Sprintf("%v", solicitud["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["Id"]), TipoEstado)

		if errTipoEstado == nil {
			*resultado = append(*resultado, map[string]interface{}{
				"Id":                    solicitud["SolicitudId"].(map[string]interface{})["Id"],
				"FechaRadicacion":       solicitud["SolicitudId"].(map[string]interface{})["FechaRadicacion"],
				"EstadoTipoSolicitudId": (*TipoEstado)["Data"].([]interface{})[0],
			})
		}
	}
}

func SolicitudGetAllSolicitudes(query string, Solicitudes []map[string]interface{}, TipoEstado *map[string]interface{}, resultado *[]interface{}, errorGetAll bool) interface{} {
	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante?limit=0"+query+"&fields=SolicitudId", &Solicitudes)

	if errSolicitud == nil {
		if Solicitudes != nil && fmt.Sprintf("%v", Solicitudes[0]) != "map[]" {
			manejoSolicitudesGetAll(Solicitudes, TipoEstado, resultado)
			return nil
		} else {
			errorGetAll = true
			// c.Data["message"] = "Error service GetAll: No data found"
			return map[string]interface{}{"Success": true, "Status": "404", "Message": "Error service GetAll: No data found", "Data": nil}
			// c.Abort("404")
		}
	} else {
		errorGetAll = true
		// c.Data["message"] = "Error service GetAll: " + errSolicitud.Error()
		return map[string]interface{}{"Success": true, "Status": "400", "Message": "Error service GetAll: " + errSolicitud.Error(), "Data": nil}
		// c.Abort("400")
	}
}

// FUNCIONES QUE SE USAN EN PUT

func iterarRespuestaSolicitudPut(RespuestaSolicitud map[string]interface{}, resDocs *[]interface{}) {
	for i := range RespuestaSolicitud["Documentos"].([]interface{}) {
		var nuevo = true

		auxDoc := []map[string]interface{}{}
		documento := map[string]interface{}{
			"IdTipoDocumento": RespuestaSolicitud["Documentos"].([]interface{})[i].(map[string]interface{})["IdTipoDocumento"],
			"nombre":          RespuestaSolicitud["Documentos"].([]interface{})[i].(map[string]interface{})["nombre"],
			"metadatos":       RespuestaSolicitud["Documentos"].([]interface{})[i].(map[string]interface{})["metadatos"],
			"descripcion":     RespuestaSolicitud["Documentos"].([]interface{})[i].(map[string]interface{})["descripcion"],
			"file":            RespuestaSolicitud["Documentos"].([]interface{})[i].(map[string]interface{})["file"],
		}
		auxDoc = append(auxDoc, documento)
		doc, errDoc := models.RegistrarDoc(auxDoc)
		if errDoc == nil {
			docTem := map[string]interface{}{
				"Nombre":        doc.(map[string]interface{})["Nombre"].(string),
				"Enlace":        doc.(map[string]interface{})["Enlace"],
				"Id":            doc.(map[string]interface{})["Id"],
				"TipoDocumento": doc.(map[string]interface{})["TipoDocumento"],
				"Activo":        doc.(map[string]interface{})["Activo"],
			}

			for index, documento := range *resDocs {
				if documento.(map[string]interface{})["TipoDocumento"].(map[string]interface{})["Id"] == RespuestaSolicitud["Documentos"].([]interface{})[i].(map[string]interface{})["IdTipoDocumento"] {
					nuevo = false
					(*resDocs)[index] = docTem
				}
			}
			if nuevo {
				*resDocs = append(*resDocs, docTem)
			}
		}
	}
}

func validarRespuestaSolicitudPut(RespuestaSolicitud map[string]interface{}, resDocs *[]interface{}, Referencia *string) {
	if len(RespuestaSolicitud["Documentos"].([]interface{})) > 0 {
		iterarRespuestaSolicitudPut(RespuestaSolicitud, resDocs)
	}

	jsonPeriodo, _ := json.Marshal(RespuestaSolicitud["Periodo"])
	jsonDocumento, _ := json.Marshal(resDocs)
	jsonProyecto, _ := json.Marshal(RespuestaSolicitud["Proyecto"])
	jsonEspacio, _ := json.Marshal(RespuestaSolicitud["EspacioAcademico"])
	jsonVehiculo, _ := json.Marshal(RespuestaSolicitud["TipoVehiculo"])
	jsonDocente, _ := json.Marshal(RespuestaSolicitud["DocenteSolicitante"])
	jsonDocentes, _ := json.Marshal(RespuestaSolicitud["DocentesInvitados"])

	*Referencia = "{\n\"Periodo\":" + fmt.Sprintf("%v", string(jsonPeriodo)) +
		",\n\"Proyecto\": " + fmt.Sprintf("%v", string(jsonProyecto)) +
		",\n\"EspacioAcademico\": " + fmt.Sprintf("%v", string(jsonEspacio)) +
		",\n\"Semestre\": " + fmt.Sprintf("%v", RespuestaSolicitud["Semestre"]) +
		",\n\"NumeroEstudiantes\": " + fmt.Sprintf("%v", RespuestaSolicitud["NumeroEstudiantes"]) +
		",\n\"NumeroGrupos\": " + fmt.Sprintf("%v", RespuestaSolicitud["NumeroGrupos"]) +
		",\n\"Duracion\": " + fmt.Sprintf("%v", RespuestaSolicitud["Duracion"]) +
		",\n\"NumeroVehiculos\": " + fmt.Sprintf("%v", RespuestaSolicitud["NumeroVehiculos"]) +
		",\n\"TipoVehiculo\": " + fmt.Sprintf("%v", string(jsonVehiculo)) +
		",\n\"FechaHoraSalida\": \"" + time_bogota.TiempoCorreccionFormato(RespuestaSolicitud["FechaHoraSalida"].(string)) + "\"" +
		",\n\"FechaHoraRegreso\": \"" + time_bogota.TiempoCorreccionFormato(RespuestaSolicitud["FechaHoraRegreso"].(string)) + "\"" +
		",\n\"Documentos\": " + fmt.Sprintf("%v", string(jsonDocumento)) +
		",\n\"DocenteSolicitante\": " + fmt.Sprintf("%v", string(jsonDocente)) +
		",\n\"DocentesInvitados\": " + fmt.Sprintf("%v", string(jsonDocentes)) + "\n}"
}

func solicitudPut(id_practica string, SolicitudPut *map[string]interface{}, Solicitud map[string]interface{}, resultado *map[string]interface{}, errorGetAll *bool, alertas *[]interface{}, alerta *models.Alert) interface{} {
	errPutEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+id_practica, "PUT", SolicitudPut, Solicitud)

	if errPutEstado == nil {
		if (*SolicitudPut)["Status"] != "400" {
			*resultado = *SolicitudPut
			return nil
		} else {
			ManejoError(alerta, alertas, fmt.Sprintf("%v", (*SolicitudPut)["Message"]), errorGetAll)
			return map[string]interface{}{"Response": *alerta}
		}
	} else {
		ManejoError(alerta, alertas, "", errorGetAll, errPutEstado)
		return map[string]interface{}{"Response": *alerta}
	}
}

func solicitudEvolucionEstadoPut(NuevoEstado map[string]interface{}, id_practica string, RespuestaSolicitud map[string]interface{}, SolicitudEvolucionEstadoPost *map[string]interface{}, SolicitudPut *map[string]interface{}, Solicitud map[string]interface{}, resultado *map[string]interface{}, errorGetAll *bool, alertas *[]interface{}, alerta *models.Alert, Referencia string, idEstado string, observacionPost *map[string]interface{}) interface{} {
	estadoId := NuevoEstado["Data"]

	id, _ := strconv.Atoi(id_practica)
	SolicitudEvolucionEstado := map[string]interface{}{
		"TerceroId": int(RespuestaSolicitud["IdTercero"].(float64)),
		"SolicitudId": map[string]interface{}{
			"Id": id,
		},
		"EstadoTipoSolicitudId": map[string]interface{}{
			"Id": int(estadoId.([]interface{})[0].(map[string]interface{})["Id"].(float64)),
		},
		"EstadoTipoSolicitudIdAnterior": map[string]interface{}{
			"Id": int(RespuestaSolicitud["EstadoTipoSolicitudIdAnterior"].(map[string]interface{})["Id"].(float64)),
		},
		"Activo":      true,
		"FechaLimite": RespuestaSolicitud["FechaRespuesta"],
	}

	errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado", "POST", SolicitudEvolucionEstadoPost, SolicitudEvolucionEstado)
	if errSolicitudEvolucionEstado == nil {
		if *SolicitudEvolucionEstadoPost != nil && fmt.Sprintf("%v", *SolicitudEvolucionEstadoPost) != "map[]" {

			Solicitud["Resultado"] = "{\"Periodo\":\"" + fmt.Sprintf("%v", string(RespuestaSolicitud["Comentario"].(string))) + "\"}"
			Solicitud["EstadoTipoSolicitudId"] = (*SolicitudEvolucionEstadoPost)["Data"].(map[string]interface{})["EstadoTipoSolicitudId"]
			Solicitud["EstadoTipoSolicitudId"].(map[string]interface{})["Activo"] = true

			// Si hay modificaciones en la información de la solicitud
			if len(Referencia) > 0 || Referencia != "" {
				Solicitud["Referencia"] = Referencia
			}

			// Si la practica es ejecutada, se da por finalizada la solicitud
			if idEstado == "23" {
				Solicitud["SolicitudFinalizada"] = true
			}

			Observacion := map[string]interface{}{
				"TerceroId": RespuestaSolicitud["IdTercero"],
				"TipoObservacionId": map[string]interface{}{
					"Id": 1,
				},
				"SolicitudId": map[string]interface{}{
					"Id": int(Solicitud["Id"].(float64)),
				},
				"Valor":  RespuestaSolicitud["Comentario"].(string),
				"Titulo": fmt.Sprintf("%v", (*SolicitudEvolucionEstadoPost)["Data"].(map[string]interface{})["Id"]),
				"Activo": true,
			}

			errObservacion := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"observacion", "POST", observacionPost, Observacion)
			if errObservacion == nil {
				return solicitudPut(id_practica, SolicitudPut, Solicitud, resultado, errorGetAll, alertas, alerta)
			}
		} else {
			ManejoError(alerta, alertas, "No data found", errorGetAll)
			return map[string]interface{}{"Response": *alerta}
		}
	} else {
		ManejoError(alerta, alertas, "", errorGetAll, errSolicitudEvolucionEstado)
		return map[string]interface{}{"Response": *alerta}
	}
	return nil
}

func solicitudNuevoEstadoPut(NuevoEstado *map[string]interface{}, id_practica string, RespuestaSolicitud map[string]interface{}, SolicitudEvolucionEstadoPost *map[string]interface{}, SolicitudPut *map[string]interface{}, Solicitud map[string]interface{}, resultado *map[string]interface{}, errorGetAll *bool, alertas *[]interface{}, alerta *models.Alert, Referencia string, idEstado string, observacionPost *map[string]interface{}, tipoSolicitud map[string]interface{}) interface{} {
	var id = fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0].(map[string]interface{})["Id"])

	errEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado_tipo_solicitud?query=EstadoId.Id:"+
		idEstado+",TipoSolicitud.Id:"+id, NuevoEstado)

	if errEstado == nil {
		return solicitudEvolucionEstadoPut(*NuevoEstado, id_practica, RespuestaSolicitud, SolicitudEvolucionEstadoPost, SolicitudPut, Solicitud, resultado, errorGetAll, alertas, alerta, Referencia, idEstado, observacionPost)
	} else {
		ManejoError(alerta, alertas, "No data found", errorGetAll)
		return map[string]interface{}{"Response": *alerta}
	}
}

func solicitudAnteriorEstadoPut(NuevoEstado *map[string]interface{}, id_practica string, RespuestaSolicitud map[string]interface{}, SolicitudEvolucionEstadoPost *map[string]interface{}, SolicitudPut *map[string]interface{}, Solicitud map[string]interface{}, resultado *map[string]interface{}, errorGetAll *bool, alertas *[]interface{}, alerta *models.Alert, Referencia string, observacionPost *map[string]interface{}, tipoSolicitud *map[string]interface{}, anteriorEstado []map[string]interface{}, anteriorEstadoPost *map[string]interface{}) interface{} {
	anteriorEstado[0]["Activo"] = false
	estasAnteriorId := fmt.Sprintf("%v", anteriorEstado[0]["Id"])

	errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado/"+estasAnteriorId, "PUT", anteriorEstadoPost, anteriorEstado[0])
	if errSolicitudEvolucionEstado == nil {

		// Búsqueda de estado relacionado con las prácticas académicas
		idEstado := fmt.Sprintf("%v", RespuestaSolicitud["Estado"].(map[string]interface{})["Id"])
		errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"tipo_solicitud?query=CodigoAbreviacion:SoPA", tipoSolicitud)
		if errTipoSolicitud == nil && fmt.Sprintf("%v", (*tipoSolicitud)["Data"].([]interface{})[0]) != "map[]" {
			return solicitudNuevoEstadoPut(NuevoEstado, id_practica, RespuestaSolicitud, SolicitudEvolucionEstadoPost, SolicitudPut, Solicitud, resultado, errorGetAll, alertas, alerta, Referencia, idEstado, observacionPost, *tipoSolicitud)
		}

	} else {
		ManejoError(alerta, alertas, "", errorGetAll, errSolicitudEvolucionEstado)
		return map[string]interface{}{"Response": *alerta}
	}
	return nil
}

func solicitudActualizarAnterioEstadoPut(NuevoEstado *map[string]interface{}, id_practica string, RespuestaSolicitud map[string]interface{}, SolicitudEvolucionEstadoPost *map[string]interface{}, SolicitudPut *map[string]interface{}, Solicitud map[string]interface{}, resultado *map[string]interface{}, errorGetAll *bool, alertas *[]interface{}, alerta *models.Alert, Referencia string, observacionPost *map[string]interface{}, tipoSolicitud *map[string]interface{}, anteriorEstado *[]map[string]interface{}, anteriorEstadoPost *map[string]interface{}) (interface{}, interface{}, bool) {
	errAntEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado?query=activo:true,solicitudId.Id:"+id_practica, anteriorEstado)
	if errAntEstado == nil {
		if *anteriorEstado != nil && fmt.Sprintf("%v", *anteriorEstado) != "map[]" {
			dataJson := solicitudAnteriorEstadoPut(NuevoEstado, id_practica, RespuestaSolicitud, SolicitudEvolucionEstadoPost, SolicitudPut, Solicitud, resultado, errorGetAll, alertas, alerta, Referencia, observacionPost, tipoSolicitud, *anteriorEstado, anteriorEstadoPost)

			if dataJson == nil {
				return nil, dataJson, true
			} else {
				return nil, dataJson, false
			}
		} else {
			ManejoError(alerta, alertas, "No data found", errorGetAll)
			return nil, map[string]interface{}{"Response": *alerta}, false
		}

	} else {
		*errorGetAll = true
		return "Error service GetAll: No data found", nil, false
	}
}

func ManejoSolicitudesPut(NuevoEstado *map[string]interface{}, id_practica string, RespuestaSolicitud map[string]interface{}, SolicitudEvolucionEstadoPost *map[string]interface{}, SolicitudPut *map[string]interface{}, Solicitud *map[string]interface{}, resultado *map[string]interface{}, errorGetAll *bool, alertas *[]interface{}, alerta *models.Alert, Referencia *string, observacionPost *map[string]interface{}, tipoSolicitud *map[string]interface{}, anteriorEstado *[]map[string]interface{}, anteriorEstadoPost *map[string]interface{}, resDocs *[]interface{}) (interface{}, interface{}, bool) {
	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+id_practica, Solicitud)
	if errSolicitud == nil {
		if *Solicitud != nil && fmt.Sprintf("%v", (*Solicitud)["Status"]) != "404" {

			var sol map[string]interface{}
			if errSol := json.Unmarshal([]byte((*Solicitud)["Referencia"].(string)), &sol); errSol == nil {
				*resDocs = sol["Documentos"].([]interface{})

				if RespuestaSolicitud["Documentos"] != nil {
					validarRespuestaSolicitudPut(RespuestaSolicitud, resDocs, Referencia)
				}

				// Actualización del anterior estado
				dataMessage, dataJson, exito := solicitudActualizarAnterioEstadoPut(NuevoEstado, id_practica, RespuestaSolicitud, SolicitudEvolucionEstadoPost, SolicitudPut, *Solicitud, resultado, errorGetAll, alertas, alerta, *Referencia, observacionPost, tipoSolicitud, anteriorEstado, anteriorEstadoPost)

				if !exito {
					if dataMessage != nil {
						return dataMessage, nil, false
						//c.Data["message"] = dataMessage
						//c.Abort("404")
					} else {
						//c.Data["json"] = dataJson
						return nil, dataJson, false
					}
				} else {
					return nil, nil, true
				}
			} else {
				*errorGetAll = true
				return "Error service GetAll: No data found", nil, false
			}

		} else {
			*errorGetAll = true
			return "Error service GetAll: No data found", nil, false
		}

	} else {
		ManejoError(alerta, alertas, "", errorGetAll, errSolicitud)
		return nil, map[string]interface{}{"Response": *alerta}, false
	}
}

// FUNCIONES QUE SE USAN EN CONSULTAR INFO SOLICITANTE

func manejoVinculacionSolicitante(resultado *map[string]interface{}, tipoVinculacion []map[string]interface{}) {
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

func SolicitudCorreoInstitucionalConsultarInfo(idTercero string, correoInstitucional *[]map[string]interface{}, jsondata *map[string]interface{}, resultado *map[string]interface{}, alertas *[]interface{}, alerta *models.Alert, errorGetAll *bool) interface{} {
	errCorreoIns := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:94,Activo:true", correoInstitucional)
	if errCorreoIns == nil && fmt.Sprintf("%v", (*correoInstitucional)[0]) != "map[]" {
		if (*correoInstitucional)[0]["Status"] != 404 {
			correoaux := (*correoInstitucional)[0]["Dato"]
			if err := json.Unmarshal([]byte(correoaux.(string)), jsondata); err != nil {
				panic(err)
			}
			(*resultado)["CorreoInstitucional"] = (*jsondata)["value"]
		} else {
			if (*correoInstitucional)[0]["Message"] == "Not found resource" {
				return nil
			} else {
				ManejoError(alerta, alertas, "No data found", errorGetAll)
				return map[string]interface{}{"Response": *alerta}
			}
		}
	} else {
		ManejoError(alerta, alertas, "No data found", errorGetAll)
		return map[string]interface{}{"Response": *alerta}
	}
	return nil
}

func SolicitudCorreoConsultarInfo(idTercero string, correoElectronico *[]map[string]interface{}, jsondata *map[string]interface{}, resultado *map[string]interface{}, alertas *[]interface{}, alerta *models.Alert, errorGetAll *bool) interface{} {
	errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:53,Activo:true", correoElectronico)
	if errCorreo == nil && fmt.Sprintf("%v", (*correoElectronico)[0]) != "map[]" {
		if (*correoElectronico)[0]["Status"] != 404 {
			correoaux := (*correoElectronico)[0]["Dato"]
			if err := json.Unmarshal([]byte(correoaux.(string)), jsondata); err != nil {
				panic(err)
			}
			(*resultado)["Correo"] = (*jsondata)["Data"]
		} else {
			if (*correoElectronico)[0]["Message"] == "Not found resource" {
				return nil
			} else {
				ManejoError(alerta, alertas, "No data found", errorGetAll)
				return map[string]interface{}{"Response": *alerta}
			}
		}
	} else {
		ManejoError(alerta, alertas, "No data found", errorGetAll)
		return map[string]interface{}{"Response": *alerta}
	}
	return nil
}

func SolicitudCorreoPersonalConsultarInfo(idTercero string, correoPersonal *[]map[string]interface{}, jsondata *map[string]interface{}, resultado *map[string]interface{}, alertas *[]interface{}, alerta *models.Alert, errorGetAll *bool) interface{} {
	errCorreoPersonal := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:253,Activo:true", correoPersonal)
	if errCorreoPersonal == nil && fmt.Sprintf("%v", (*correoPersonal)[0]) != "map[]" {
		if (*correoPersonal)[0]["Status"] != 404 {
			correoaux := (*correoPersonal)[0]["Dato"]
			if err := json.Unmarshal([]byte(correoaux.(string)), jsondata); err != nil {
				panic(err)
			}
			(*resultado)["CorreoPersonal"] = (*jsondata)["Data"]
		} else {
			if (*correoPersonal)[0]["Message"] == "Not found resource" {
				return nil
			} else {
				ManejoError(alerta, alertas, "No data found", errorGetAll)
				return map[string]interface{}{"Response": *alerta}
			}
		}
	} else {
		ManejoError(alerta, alertas, "No data found", errorGetAll)
		return map[string]interface{}{"Response": *alerta}
	}
	return nil
}

func SolicitudTelefonoConsultarInfo(idTercero string, telefono *[]map[string]interface{}, jsondata *map[string]interface{}, resultado *map[string]interface{}, alertas *[]interface{}, alerta *models.Alert, errorGetAll *bool) interface{} {
	errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:51,Activo:true", telefono)
	if errTelefono == nil && fmt.Sprintf("%v", (*telefono)[0]) != "map[]" {
		if (*telefono)[0]["Status"] != 404 {
			telefonoaux := (*telefono)[0]["Dato"]

			if err := json.Unmarshal([]byte(telefonoaux.(string)), jsondata); err != nil {
				(*resultado)["Telefono"] = (*telefono)[0]["Dato"]
			} else {
				(*resultado)["Telefono"] = (*jsondata)["principal"]
			}
		} else {
			if (*telefono)[0]["Message"] == "Not found resource" {
				return nil
			} else {
				ManejoError(alerta, alertas, "No data found", errorGetAll)
				return map[string]interface{}{"Response": *alerta}
			}
		}
	} else {
		ManejoError(alerta, alertas, "No data found", errorGetAll)
		return map[string]interface{}{"Response": *alerta}
	}
	return nil
}

func SolicitudCelularConsultarInfo(idTercero string, celular *[]map[string]interface{}, jsondata *map[string]interface{}, resultado *map[string]interface{}, alertas *[]interface{}, alerta *models.Alert, errorGetAll *bool) interface{} {
	errCelular := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:52,Activo:true", celular)
	if errCelular == nil && fmt.Sprintf("%v", (*celular)[0]) != "map[]" {
		if (*celular)[0]["Status"] != 404 {
			celularaux := (*celular)[0]["Dato"]

			if err := json.Unmarshal([]byte(celularaux.(string)), &jsondata); err != nil {
				(*resultado)["Celular"] = (*celular)[0]["Dato"]
			} else {
				(*resultado)["Celular"] = (*jsondata)["principal"]
			}
		} else {
			if (*celular)[0]["Message"] == "Not found resource" {
				return nil
			} else {
				ManejoError(alerta, alertas, "No data found", errorGetAll)
				return map[string]interface{}{"Response": *alerta}
			}
		}
	} else {
		ManejoError(alerta, alertas, "No data found", errorGetAll)
		return map[string]interface{}{"Response": *alerta}
	}
	return nil
}

func SolicitudTipoVinculacionConsultarInfo(idTercero string, tipoVinculacion *[]map[string]interface{}, resultado *map[string]interface{}, alertas *[]interface{}, alerta *models.Alert, errorGetAll *bool) interface{} {
	errVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"vinculacion?query=TerceroPrincipalId:"+fmt.Sprintf("%v", idTercero)+",Activo:true&limit=0", tipoVinculacion)
	if errVinculacion == nil && fmt.Sprintf("%v", (*tipoVinculacion)[0]) != "map[]" {
		if (*tipoVinculacion)[0]["Status"] != 404 {
			manejoVinculacionSolicitante(resultado, *tipoVinculacion)
		} else {
			if (*tipoVinculacion)[0]["Message"] == "Not found resource" {
				return nil
			} else {
				ManejoError(alerta, alertas, "No data found", errorGetAll)
				return map[string]interface{}{"Response": *alerta}
			}
		}
	} else {
		ManejoError(alerta, alertas, "No data found", errorGetAll)
		return map[string]interface{}{"Response": *alerta}
	}
	return nil
}

// FUNCIONES QUE SE USAN EN CONSULTAR INFO COLABORADOR

func solicitudVinculacionInfoColaborador(tv map[string]interface{}, vinculacion *map[string]interface{}, resultado *map[string]interface{}) {
	errVinculacion := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro?query=Id:"+fmt.Sprintf("%v", tv["TipoVinculacionId"])+",Activo:true&limit=0", vinculacion)
	if errVinculacion == nil && fmt.Sprintf("%v", (*vinculacion)["Data"]) != "[map[]]" {
		if (*vinculacion)["Status"] != 404 {
			(*resultado)["TipoVinculacionId"] = (*vinculacion)["Data"].([]interface{})[0]
		}
	}
}

func solicitudCorreoInstitucionalInfoColaborador(idTercero interface{}, correoInstitucional *[]map[string]interface{}, jsondata *map[string]interface{}, resultado *map[string]interface{}, alertas *[]interface{}, alerta *models.Alert, errorGetAll *bool) interface{} {
	errCorreoIns := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:94,Activo:true", correoInstitucional)
	if errCorreoIns == nil && fmt.Sprintf("%v", (*correoInstitucional)[0]) != "map[]" {
		if (*correoInstitucional)[0]["Status"] != 404 {
			correoaux := (*correoInstitucional)[0]["Dato"]
			if err := json.Unmarshal([]byte(correoaux.(string)), jsondata); err != nil {
				panic(err)
			}
			(*resultado)["CorreoInstitucional"] = (*jsondata)["value"]
			return nil
		} else {
			if (*correoInstitucional)[0]["Message"] == "Not found resource" {
				return nil
			} else {
				ManejoError(alerta, alertas, "No data found", errorGetAll)
				return map[string]interface{}{"Response": *alerta}
			}
		}
	} else {
		ManejoError(alerta, alertas, "No data found", errorGetAll)
		return map[string]interface{}{"Response": *alerta}
	}
}

func solicitudCorreoInfoColaborador(idTercero interface{}, correoElectronico *[]map[string]interface{}, jsondata *map[string]interface{}, resultado *map[string]interface{}, alertas *[]interface{}, alerta *models.Alert, errorGetAll *bool) interface{} {
	errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:53,Activo:true", correoElectronico)
	if errCorreo == nil && fmt.Sprintf("%v", (*correoElectronico)[0]) != "map[]" {
		if (*correoElectronico)[0]["Status"] != 404 {
			correoaux := (*correoElectronico)[0]["Dato"]
			if err := json.Unmarshal([]byte(correoaux.(string)), jsondata); err != nil {
				panic(err)
			}
			(*resultado)["Correo"] = (*jsondata)["Data"]
			return nil
		} else {
			if (*correoElectronico)[0]["Message"] == "Not found resource" {
				return nil
			} else {
				ManejoError(alerta, alertas, "No data found", errorGetAll)
				return map[string]interface{}{"Response": *alerta}
			}
		}
	} else {
		ManejoError(alerta, alertas, "No data found", errorGetAll)
		return map[string]interface{}{"Response": *alerta}
	}
}

func solicitudCorreoPersonalInfoColaborador(idTercero interface{}, correoPersonal *[]map[string]interface{}, jsondata *map[string]interface{}, resultado *map[string]interface{}, alertas *[]interface{}, alerta *models.Alert, errorGetAll *bool) interface{} {
	errCorreoPersonal := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:253,Activo:true", correoPersonal)
	if errCorreoPersonal == nil && fmt.Sprintf("%v", (*correoPersonal)[0]) != "map[]" {
		if (*correoPersonal)[0]["Status"] != 404 {
			correoaux := (*correoPersonal)[0]["Dato"]
			if err := json.Unmarshal([]byte(correoaux.(string)), jsondata); err != nil {
				panic(err)
			}
			(*resultado)["CorreoPersonal"] = (*jsondata)["Data"]
			return nil
		} else {
			if (*correoPersonal)[0]["Message"] == "Not found resource" {
				return nil
			} else {
				ManejoError(alerta, alertas, "No data found", errorGetAll)
				return map[string]interface{}{"Response": *alerta}
			}
		}
	} else {
		ManejoError(alerta, alertas, "No data found", errorGetAll)
		return map[string]interface{}{"Response": *alerta}
	}
}

func solicitudTelefonoInfoColaborador(idTercero interface{}, telefono *[]map[string]interface{}, jsondata *map[string]interface{}, resultado *map[string]interface{}, alertas *[]interface{}, alerta *models.Alert, errorGetAll *bool) interface{} {
	errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:51,Activo:true", telefono)
	if errTelefono == nil && fmt.Sprintf("%v", (*telefono)[0]) != "map[]" {
		if (*telefono)[0]["Status"] != 404 {
			telefonoaux := (*telefono)[0]["Dato"]

			if err := json.Unmarshal([]byte(telefonoaux.(string)), jsondata); err != nil {
				(*resultado)["Telefono"] = (*telefono)[0]["Dato"]
			} else {
				(*resultado)["Telefono"] = (*jsondata)["principal"]
			}
			return nil
		} else {
			if (*telefono)[0]["Message"] == "Not found resource" {
				return nil
			} else {
				ManejoError(alerta, alertas, "No data found", errorGetAll)
				return map[string]interface{}{"Response": *alerta}
			}
		}
	} else {
		ManejoError(alerta, alertas, "No data found", errorGetAll)
		return map[string]interface{}{"Response": *alerta}
	}
}

func solicitudCelularInfoColaborador(idTercero interface{}, celular *[]map[string]interface{}, jsondata *map[string]interface{}, resultado *map[string]interface{}, alertas *[]interface{}, alerta *models.Alert, errorGetAll *bool) interface{} {
	errCelular := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:52,Activo:true", celular)
	if errCelular == nil && fmt.Sprintf("%v", (*celular)[0]) != "map[]" {
		if (*celular)[0]["Status"] != 404 {
			celularaux := (*celular)[0]["Dato"]

			if err := json.Unmarshal([]byte(celularaux.(string)), jsondata); err != nil {
				(*resultado)["Celular"] = (*celular)[0]["Dato"]
			} else {
				(*resultado)["Celular"] = (*jsondata)["principal"]
			}
			return nil
		} else {
			if (*celular)[0]["Message"] == "Not found resource" {
				return nil
			} else {
				ManejoError(alerta, alertas, "No data found", errorGetAll)
				return map[string]interface{}{"Response": *alerta}
			}
		}
	} else {
		ManejoError(alerta, alertas, "No data found", errorGetAll)
		return map[string]interface{}{"Response": *alerta}
	}
}

func iteracionTipoVinculacionInfoColaborador(tipoVinculacion []map[string]interface{}, idTercero interface{}, jsondata *map[string]interface{}, resultado *map[string]interface{}, alertas *[]interface{}, alerta *models.Alert, errorGetAll *bool, correoInstitucional *[]map[string]interface{}, correoElectronico *[]map[string]interface{}, correoPersonal *[]map[string]interface{}, telefono *[]map[string]interface{}, celular *[]map[string]interface{}, persona []map[string]interface{}) interface{} {
	var data interface{}
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

			solicitudVinculacionInfoColaborador(tv, &vinculacion, resultado)

			// Correo institucional --> 94
			data = solicitudCorreoInstitucionalInfoColaborador(idTercero, correoInstitucional, jsondata, resultado, alertas, alerta, errorGetAll)

			// Correo --> 53
			data = solicitudCorreoInfoColaborador(idTercero, correoElectronico, jsondata, resultado, alertas, alerta, errorGetAll)

			// Correo personal --> 253
			data = solicitudCorreoPersonalInfoColaborador(idTercero, correoPersonal, jsondata, resultado, alertas, alerta, errorGetAll)

			// Teléfono --> 51
			data = solicitudTelefonoInfoColaborador(idTercero, telefono, jsondata, resultado, alertas, alerta, errorGetAll)

			// Celular --> 52
			data = solicitudCelularInfoColaborador(idTercero, celular, jsondata, resultado, alertas, alerta, errorGetAll)

			(*resultado)["Nombre"] = persona[0]["TerceroId"].(map[string]interface{})["NombreCompleto"]
			(*resultado)["Id"] = idTercero
			(*resultado)["PuedeBorrar"] = true
			break
		} else {
			logs.Error("No es docente")
			return map[string]interface{}{"Success": false, "Status": "404", "Message": "Data not found", "Data": nil}
		}
	}
	return data
}

func SolicitudTipoVinculacionInfoColaborador(tipoVinculacion *[]map[string]interface{}, idTercero interface{}, jsondata *map[string]interface{}, resultado *map[string]interface{}, alertas *[]interface{}, alerta *models.Alert, errorGetAll *bool, correoInstitucional *[]map[string]interface{}, correoElectronico *[]map[string]interface{}, correoPersonal *[]map[string]interface{}, telefono *[]map[string]interface{}, celular *[]map[string]interface{}, persona []map[string]interface{}) interface{} {
	errVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"vinculacion?query=TerceroPrincipalId:"+fmt.Sprintf("%v", idTercero)+",Activo:true&limit=0", tipoVinculacion)
	if errVinculacion == nil && fmt.Sprintf("%v", (*tipoVinculacion)[0]) != "map[]" {
		if (*tipoVinculacion)[0]["Status"] != 404 {
			return iteracionTipoVinculacionInfoColaborador(*tipoVinculacion, idTercero, jsondata, resultado, alertas, alerta, errorGetAll, correoInstitucional, correoElectronico, correoPersonal, telefono, celular, persona)
		} else {
			if (*tipoVinculacion)[0]["Message"] == "Not found resource" {
				return nil
			} else {
				ManejoError(alerta, alertas, "No data found", errorGetAll)
				return map[string]interface{}{"Success": false, "Status": "404", "Message": "Data not found", "Data": nil}
			}
		}
	} else {
		ManejoError(alerta, alertas, "No data found", errorGetAll)
		return map[string]interface{}{"Success": false, "Status": "404", "Message": "Data not found", "Data": nil}
	}
}

// FUNCIONES QUE SE USAN EN CONSULTAR PARAMETROS

func manejoProyectosParametros(resultado *map[string]interface{}, getProyecto *[]map[string]interface{}, proyectos []map[string]interface{}) {
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

func manejoEstadosParametros(resultado *map[string]interface{}, tipoEstados map[string]interface{}, estados []interface{}) {
	if tipoEstados["Status"] != "404" {
		for _, estado := range tipoEstados["Data"].([]interface{}) {
			estados = append(estados, estado.(map[string]interface{})["EstadoId"])
			(*resultado)["estados"] = estados
		}
	}
}

func SolicitudPeriodoParametros(periodos *map[string]interface{}, resultado *map[string]interface{}) interface{} {
	errPeriodo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo?query=CodigoAbreviacion:PA,Activo:true&limit=1&sortby=Id&order=desc", periodos)
	if errPeriodo == nil && fmt.Sprintf("%v", (*periodos)["Data"]) != "[map[]]" {
		if (*periodos)["Status"] != "404" {
			(*resultado)["periodos"] = (*periodos)["Data"]
		}
		return nil
	} else {
		(*resultado)["periodos"] = nil
		logs.Error(*periodos)
		return errPeriodo
	}
}

func SolicitudProyectoParametros(getProyecto *[]map[string]interface{}, resultado *map[string]interface{}, proyectos []map[string]interface{}) interface{} {
	errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion/?query=Activo:true,Oferta:true&limit=0", getProyecto)
	if errProyecto == nil {
		manejoProyectosParametros(resultado, getProyecto, proyectos)
		return nil
	} else {
		(*resultado)["proyectos"] = nil
		logs.Error(*getProyecto)
		return errProyecto
	}
}

func SolicitudVehiculoParametros(vehiculos *map[string]interface{}, resultado *map[string]interface{}, getProyecto []map[string]interface{}) interface{} {
	errVehiculo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro/?query=tipo_parametro_id:38&sortby=numero_orden&order=asc&limit=0", vehiculos)
	if errVehiculo == nil && fmt.Sprintf("%v", (*vehiculos)["Data"]) != "[map[]]" {
		if (*vehiculos)["Status"] != "404" {
			(*resultado)["vehiculos"] = (*vehiculos)["Data"]
		}
		return nil
	} else {
		(*resultado)["proyectos"] = nil
		logs.Error(getProyecto)
		return errVehiculo
	}
}

func SolicitudTipoParametros(tipoSolicitud *map[string]interface{}, tipoEstados *map[string]interface{}, resultado *map[string]interface{}, estados []interface{}) interface{} {
	errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"tipo_solicitud?query=CodigoAbreviacion:SoPA", tipoSolicitud)
	if errTipoSolicitud == nil && fmt.Sprintf("%v", (*tipoSolicitud)["Data"].([]interface{})[0]) != "map[]" {
		var id = fmt.Sprintf("%v", (*tipoSolicitud)["Data"].([]interface{})[0].(map[string]interface{})["Id"])

		errTipoEstados := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado_tipo_solicitud?query=TipoSolicitud.Id:"+id, tipoEstados)
		if errTipoEstados == nil && fmt.Sprintf("%v", (*tipoEstados)["Data"]) != "[map[]]" {
			manejoEstadosParametros(resultado, *tipoEstados, estados)
		} else {
			logs.Error(*tipoEstados)
			return errTipoEstados
		}
	}
	return nil
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

// FUNCIONES QUE SE USAN EN ENVIAR INVITACIONES

func solicitudCorreoPostEnviarInvitaciones(CorreoPost *map[string]interface{}, correo map[string]interface{}, errorGetAll *bool, alertas *[]interface{}, alerta *models.Alert) interface{} {
	errEnvioCorreos := request.SendJson("http://"+beego.AppConfig.String("GOOGLE_MID")+"notificacion", "POST", CorreoPost, correo)
	if errEnvioCorreos == nil {
		if *CorreoPost == nil || fmt.Sprintf("%v", *CorreoPost) == "400" {
			ManejoError(alerta, alertas, "No data found", errorGetAll)
			return map[string]interface{}{"Response": *alerta}
		}
	} else {
		ManejoError(alerta, alertas, "No data found", errorGetAll)
		return map[string]interface{}{"Response": *alerta}
	}
	return nil
}

func ValidarEstadoEnviarInvitaciones(idEstado string, ReferenciaJson map[string]interface{}, CorreoPost *map[string]interface{}, errorGetAll *bool, alertas *[]interface{}, alerta *models.Alert) interface{} {
	if idEstado == "39" {

		// TO DO: Consulta de correos electronicos de los estudiantes inscritos en el espacio académico
		correoEstudiantes := []interface{}{
			"correo1@correo.com", "correo2@correo.com",
		}
		nombreEstudiantes := []interface{}{
			"Nombre 1", "Nombre 2",
		}

		var data interface{}
		for index, correo := range correoEstudiantes {
			correo := map[string]interface{}{
				"to":           []interface{}{correo},
				"cc":           []interface{}{},
				"bcc":          []interface{}{},
				"subject":      "Invitación a práctica académica",
				"templateName": "invitacion_practica_academica.html",
				"templateData": map[string]interface{}{
					"Fecha":            strings.Replace(time_bogota.TiempoBogotaFormato()[:16], "T", " ", 1),
					"FechaInicio":      strings.Replace(time_bogota.TiempoCorreccionFormato(ReferenciaJson["FechaHoraSalida"].(string))[:16], "t", " ", 1),
					"FechaFin":         strings.Replace(time_bogota.TiempoCorreccionFormato(ReferenciaJson["FechaHoraRegreso"].(string))[:16], "t", " ", 1),
					"EspacioAcademico": ReferenciaJson["EspacioAcademico"].(map[string]interface{})["Nombre"],
					"NombreEstudiante": nombreEstudiantes[index],
					"NombreDocente":    ReferenciaJson["DocenteSolicitante"].(map[string]interface{})["Nombre"],
				},
			}

			fmt.Println("http://" + beego.AppConfig.String("GOOGLE_MID") + "notificacion")

			data = solicitudCorreoPostEnviarInvitaciones(CorreoPost, correo, errorGetAll, alertas, alerta)

			if data != nil {
				return data
			}
		}
		return data
	}
	return nil
}

// FUNCIONES QUE SE USAN EN VARIOS ENDPOINTS

func ManejoError(alerta *models.Alert, alertas *[]interface{}, mensaje string, errorGetAll *bool, err ...error) {
	var msj string
	if len(err) > 0 && err[0] != nil {
		msj = mensaje + err[0].Error()
	} else {
		msj = mensaje
	}
	*errorGetAll = true
	*alertas = append(*alertas, msj)
	(*alerta).Body = *alertas
	(*alerta).Type = "error"
	(*alerta).Code = "400"
}

func ManejoExito(alertas *[]interface{}, alerta *models.Alert, resultado map[string]interface{}) {
	*alertas = append(*alertas, resultado)
	(*alerta).Body = *alertas
	(*alerta).Code = "200"
	(*alerta).Type = "OK"
}
