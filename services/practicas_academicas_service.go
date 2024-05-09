package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_practicas_academicas_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
	"github.com/udistrital/utils_oas/time_bogota"
	"golang.org/x/sync/errgroup"
)

func CrearPracticaAcademica(data []byte) requestresponse.APIResponse {
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

	if err := json.Unmarshal(data, &solicitud); err == nil {
		for i := range solicitud["Documentos"].([]interface{}) {
			auxDoc := []map[string]interface{}{}
			documento := map[string]interface{}{
				"IdTipoDocumento": solicitud["Documentos"].([]interface{})[i].(map[string]interface{})["IdTipoDocumento"],
				"nombre":          solicitud["Documentos"].([]interface{})[i].(map[string]interface{})["nombre"],
				"metadatos":       solicitud["Documentos"].([]interface{})[i].(map[string]interface{})["metadatos"],
				"descripcion":     solicitud["Documentos"].([]interface{})[i].(map[string]interface{})["descripcion"],
				"file":            solicitud["Documentos"].([]interface{})[i].(map[string]interface{})["file"],
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

				resDocs = append(resDocs, docTem)
			}
		}

		jsonPeriodo, _ := json.Marshal(solicitud["Periodo"])
		jsonDocumento, _ := json.Marshal(resDocs)
		jsonProyecto, _ := json.Marshal(solicitud["Proyecto"])
		jsonEspacio, _ := json.Marshal(solicitud["EspacioAcademico"])
		jsonVehiculo, _ := json.Marshal(solicitud["TipoVehiculo"])
		jsonDocente, _ := json.Marshal(solicitud["DocenteSolicitante"])
		jsonDocentes, _ := json.Marshal(solicitud["DocentesInvitados"])

		Referencia = "{\n\"Periodo\":" + fmt.Sprintf("%v", string(jsonPeriodo)) +
			",\n\"Proyecto\": " + fmt.Sprintf("%v", string(jsonProyecto)) +
			",\n\"EspacioAcademico\": " + fmt.Sprintf("%v", string(jsonEspacio)) +
			",\n\"Semestre\": " + fmt.Sprintf("%v", solicitud["Semestre"]) +
			",\n\"NumeroEstudiantes\": " + fmt.Sprintf("%v", solicitud["NumeroEstudiantes"]) +
			",\n\"NumeroGrupos\": " + fmt.Sprintf("%v", solicitud["NumeroGrupos"]) +
			",\n\"Duracion\": " + fmt.Sprintf("%v", solicitud["Duracion"]) +
			",\n\"NumeroVehiculos\": " + fmt.Sprintf("%v", solicitud["NumeroVehiculos"]) +
			",\n\"TipoVehiculo\": " + fmt.Sprintf("%v", string(jsonVehiculo)) +
			",\n\"FechaHoraSalida\": \"" + time_bogota.TiempoCorreccionFormato(solicitud["FechaHoraSalida"].(string)) + "\"" +
			",\n\"FechaHoraRegreso\": \"" + time_bogota.TiempoCorreccionFormato(solicitud["FechaHoraRegreso"].(string)) + "\"" +
			",\n\"Documentos\": " + fmt.Sprintf("%v", string(jsonDocumento)) +
			",\n\"DocenteSolicitante\": " + fmt.Sprintf("%v", string(jsonDocente)) +
			",\n\"DocentesInvitados\": " + fmt.Sprintf("%v", string(jsonDocentes)) + "\n}"

		IdEstadoTipoSolicitud = 34

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
				resultado["Solicitud"] = SolicitudPost["Data"]
				IdSolicitud := SolicitudPost["Data"].(map[string]interface{})["Id"]

				//POST tabla solicitante
				Solicitante := map[string]interface{}{
					"TerceroId": solicitud["SolicitanteId"],
					"SolicitudId": map[string]interface{}{
						"Id": IdSolicitud,
					},
					"Activo": true,
				}

				errSolicitante := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante", "POST", &SolicitantePost, Solicitante)
				if errSolicitante == nil && fmt.Sprintf("%v", SolicitantePost["Status"]) != "400" {
					if SolicitantePost != nil && fmt.Sprintf("%v", SolicitantePost) != "map[]" {
						//POST a la tabla solicitud_evolucion estado
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

						errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado", "POST", &SolicitudEvolucionEstadoPost, SolicitudEvolucionEstado)
						if errSolicitudEvolucionEstado == nil {
							if SolicitudEvolucionEstadoPost != nil && fmt.Sprintf("%v", SolicitudEvolucionEstadoPost) != "map[]" {

								resultado["Solicitante"] = SolicitantePost["Data"]
							} else {
								errorGetAll = true
								alertas = append(alertas, "No data found")
								alerta.Code = 404
								alerta.Type = "error"
								alerta.Body = alertas
							}
						} else {
							var resultado2 map[string]interface{}
							request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+fmt.Sprintf("%v", IdSolicitud), "DELETE", &resultado2, nil)
							request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante/"+fmt.Sprintf("%v", SolicitantePost["Id"]), "DELETE", &resultado2, nil)
							errorGetAll = true
							alertas = append(alertas, errSolicitante.Error())
							alerta.Code = 400
							alerta.Type = "error"
							alerta.Body = alertas
						}
					} else {
						errorGetAll = true
						alertas = append(alertas, "No data found")
						alerta.Code = 404
						alerta.Type = "error"
						alerta.Body = alertas
					}
				} else {
					//Se elimina el registro de solicitud si no se puede hacer el POST a la tabla solicitante
					var resultado2 map[string]interface{}
					request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+fmt.Sprintf("%v", IdSolicitud), "DELETE", &resultado2, nil)
					errorGetAll = true
					alertas = append(alertas, errSolicitante.Error())
					alerta.Code = 400
					alerta.Type = "error"
					alerta.Body = alertas
				}
			} else {
				errorGetAll = true
				alertas = append(alertas, "No data found")
				alerta.Code = 404
				alerta.Type = "error"
				// alerta.Body = alertas
				alerta.Body = SolicitudPracticas
			}
		} else {
			errorGetAll = true
			alertas = append(alertas, errSolicitud.Error())
			alerta.Code = 400
			alerta.Type = "error"
			alerta.Body = alertas
		}

	} else {
		errorGetAll = true
		alertas = append(alertas, err.Error())
		alerta.Code = 400
		alerta.Type = "error"
		alerta.Body = alertas
	}

	if !errorGetAll {
		alertas = append(alertas, resultado)
		alerta.Code = 200
		alerta.Type = "OK"
		alerta.Body = alertas
		return requestresponse.APIResponseDTO(true, alerta.Code, alerta.Body)
	}
	return requestresponse.APIResponseDTO(false, alerta.Code, nil, alerta.Body)
}

func GetPracticaAcademica(id_practica string) requestresponse.APIResponse {
	var Solicitudes []map[string]interface{}
	var tipoSolicitud map[string]interface{}
	var Estados []map[string]interface{}
	var Comentario []map[string]interface{}
	resultado := make(map[string]interface{})
	var success bool = true
	var message string = ""
	var statusCode int = 200

	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante?query=SolicitudId.Id:"+id_practica, &Solicitudes)
	if errSolicitud == nil {
		if Solicitudes != nil && fmt.Sprintf("%v", Solicitudes[0]) != "map[]" {
			Referencia := Solicitudes[0]["SolicitudId"].(map[string]interface{})["Referencia"].(string)
			fechaRadicado := Solicitudes[0]["SolicitudId"].(map[string]interface{})["FechaRadicacion"].(string)
			var ReferenciaJson map[string]interface{}
			if err := json.Unmarshal([]byte(Referencia), &ReferenciaJson); err == nil {
				ReferenciaJson["Id"] = id_practica
				resultado = ReferenciaJson
				resultado["FechaRadicado"] = fechaRadicado
			}

			idEstado := fmt.Sprintf("%v", Solicitudes[0]["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["Id"].(float64))

			errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado_tipo_solicitud?query=Activo:true,Id:"+idEstado, &tipoSolicitud)
			if errTipoSolicitud == nil {
				if tipoSolicitud != nil && fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0]) != "map[]" {
					resultado["EstadoTipoSolicitudId"] = tipoSolicitud["Data"].([]interface{})[0]
				}
			}

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
					resultado["Estados"] = Estados
				}
			}

		} else {
			success = false
			message = "Error service GetAll: No data found"
			statusCode = 404
		}
	} else {
		success = false
		message = "Error service GetAll: " + errSolicitud.Error()
		statusCode = 400
	}

	if success {
		return requestresponse.APIResponseDTO(success, statusCode, resultado)
	} else {
		return requestresponse.APIResponseDTO(success, statusCode, nil, message)
	}
}

func GetAllPracticasAcademicas(query string, fields string) requestresponse.APIResponse {
	var Solicitudes []map[string]interface{}
	resultado := []interface{}{}
	var success bool = true
	var message interface{} = ""
	var statusCode int = 200
	wge := new(errgroup.Group)
var mutex sync.Mutex // Mutex para proteger el acceso a resultados


	fmt.Println("Antes de la solicitud")
	errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante?limit=0"+query+"&fields=SolicitudId", &Solicitudes)
	fmt.Println("Despues de la solicitud")
	if errSolicitud == nil {
		fmt.Println("Opcion 1")
		if Solicitudes != nil && fmt.Sprintf("%v", Solicitudes[0]) != "[map[]]" {
			fmt.Println("Opcion 2")
			wge.SetLimit(10)
			for _, solicitud := range Solicitudes {
				solicitud := solicitud
				wge.Go(func () error{
					var TipoEstado map[string]interface{}
					errTipoEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado_tipo_solicitud?query=Id:"+fmt.Sprintf("%v", solicitud["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["Id"]), &TipoEstado)
	
					if errTipoEstado == nil && fmt.Sprintf("%v",TipoEstado) != "map[]" {
						fmt.Println("Opcion 3")
						auxResultado := map[string]interface{}{
							"Id":                    solicitud["SolicitudId"].(map[string]interface{})["Id"],
							"FechaRadicacion":       solicitud["SolicitudId"].(map[string]interface{})["FechaRadicacion"],
							"EstadoTipoSolicitudId": TipoEstado["Data"].([]interface{})[0],
						}
						mutex.Lock()
						resultado = append(resultado, auxResultado)
						mutex.Unlock()
					}
					return nil
				})
			}
			//Si existe error, se realiza
			if err := wge.Wait(); err != nil {
				return requestresponse.APIResponseDTO(false, 400, err)
			}

		} else {
			fmt.Println("Opcion 4")
			success = false
			statusCode = 404
			message = "Error service GetAll: No data founError service GetAll: No data found"
		}
	} else {
		fmt.Println("Opcion 5")
		success = false
		statusCode = 400
		message = errSolicitud.Error()
	}

	if success {
		fmt.Println("Opcion 6")
		return requestresponse.APIResponseDTO(success, statusCode, resultado)
	} else {
		fmt.Println("Opcion 7")
		return requestresponse.APIResponseDTO(success, statusCode, nil, message)
	}
}

func ActualizarPracticaAcademica(id_practica string, data []byte) requestresponse.APIResponse {
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
	var success bool = true
	var message interface{} = ""
	var statusCode int = 200

	if err := json.Unmarshal(data, &RespuestaSolicitud); err == nil {

		// Consulta de información de la solicitud
		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+id_practica, &Solicitud)
		if errSolicitud == nil {
			if Solicitud != nil && fmt.Sprintf("%v", Solicitud["Status"]) != "404" {

				var sol map[string]interface{}
				if errSol := json.Unmarshal([]byte(Solicitud["Referencia"].(string)), &sol); errSol == nil {
					resDocs = sol["Documentos"].([]interface{})

					if RespuestaSolicitud["Documentos"] != nil {
						if len(RespuestaSolicitud["Documentos"].([]interface{})) > 0 {
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

									for index, documento := range resDocs {
										if documento.(map[string]interface{})["TipoDocumento"].(map[string]interface{})["Id"] == RespuestaSolicitud["Documentos"].([]interface{})[i].(map[string]interface{})["IdTipoDocumento"] {
											nuevo = false
											resDocs[index] = docTem
										}
									}
									if nuevo {
										resDocs = append(resDocs, docTem)
									}
								}
							}
						}

						jsonPeriodo, _ := json.Marshal(RespuestaSolicitud["Periodo"])
						jsonDocumento, _ := json.Marshal(resDocs)
						jsonProyecto, _ := json.Marshal(RespuestaSolicitud["Proyecto"])
						jsonEspacio, _ := json.Marshal(RespuestaSolicitud["EspacioAcademico"])
						jsonVehiculo, _ := json.Marshal(RespuestaSolicitud["TipoVehiculo"])
						jsonDocente, _ := json.Marshal(RespuestaSolicitud["DocenteSolicitante"])
						jsonDocentes, _ := json.Marshal(RespuestaSolicitud["DocentesInvitados"])

						Referencia = "{\n\"Periodo\":" + fmt.Sprintf("%v", string(jsonPeriodo)) +
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

					// Actualización del anterior estado
					errAntEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado?query=activo:true,solicitudId.Id:"+id_practica, &anteriorEstado)
					if errAntEstado == nil {
						if anteriorEstado != nil && fmt.Sprintf("%v", anteriorEstado) != "map[]" {

							anteriorEstado[0]["Activo"] = false
							estasAnteriorId := fmt.Sprintf("%v", anteriorEstado[0]["Id"])

							errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado/"+estasAnteriorId, "PUT", &anteriorEstadoPost, anteriorEstado[0])
							if errSolicitudEvolucionEstado == nil {

								// Búsqueda de estado relacionado con las prácticas académicas
								idEstado := fmt.Sprintf("%v", RespuestaSolicitud["Estado"].(map[string]interface{})["Id"])
								errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"tipo_solicitud?query=CodigoAbreviacion:SoPA", &tipoSolicitud)
								if errTipoSolicitud == nil && fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0]) != "map[]" {
									var id = fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0].(map[string]interface{})["Id"])

									errEstado := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado_tipo_solicitud?query=EstadoId.Id:"+
										idEstado+",TipoSolicitud.Id:"+id, &NuevoEstado)

									if errEstado == nil {

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

										errSolicitudEvolucionEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud_evolucion_estado", "POST", &SolicitudEvolucionEstadoPost, SolicitudEvolucionEstado)
										if errSolicitudEvolucionEstado == nil {
											if SolicitudEvolucionEstadoPost != nil && fmt.Sprintf("%v", SolicitudEvolucionEstadoPost) != "map[]" {

												Solicitud["Resultado"] = "{\"Periodo\":\"" + fmt.Sprintf("%v", string(RespuestaSolicitud["Comentario"].(string))) + "\"}"
												Solicitud["EstadoTipoSolicitudId"] = SolicitudEvolucionEstadoPost["Data"].(map[string]interface{})["EstadoTipoSolicitudId"]
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
													"Titulo": fmt.Sprintf("%v", SolicitudEvolucionEstadoPost["Data"].(map[string]interface{})["Id"]),
													"Activo": true,
												}

												errObservacion := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"observacion", "POST", &observacionPost, Observacion)
												if errObservacion == nil {

													errPutEstado := request.SendJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitud/"+id_practica, "PUT", &SolicitudPut, Solicitud)

													if errPutEstado == nil {
														if SolicitudPut["Status"] != "400" {
															resultado = SolicitudPut
														} else {
															success = false
															message = SolicitudPut["Message"]
															statusCode = 400
														}
													} else {
														success = false
														message = errPutEstado.Error()
														statusCode = 400
													}
												}
											} else {
												success = false
												message = "No data found"
												statusCode = 404
											}
										} else {
											success = false
											message = errSolicitudEvolucionEstado.Error()
											statusCode = 400
										}
									} else {
										success = false
										message = "No data found"
										statusCode = 404
									}
								}

							} else {
								success = false
								message = errSolicitudEvolucionEstado.Error()
								statusCode = 400
							}

						} else {
							success = false
							message = "No data found"
							statusCode = 404
						}

					} else {
						success = false
						message = errAntEstado.Error()
						statusCode = 400
					}

				} else {
					success = false
					message = errSol.Error()
					statusCode = 400
				}

			} else {
				success = false
				message = "No data found"
				statusCode = 404
			}

		} else {
			success = false
			message = errSolicitud.Error()
			statusCode = 400
		}

	} else {
		success = false
		message = err.Error()
		statusCode = 400
	}

	if success {
		return requestresponse.APIResponseDTO(success, statusCode, resultado)
	} else {
		return requestresponse.APIResponseDTO(success, statusCode, nil, message)
	}
}

func ConsultarInfoSolicitante(idSolicitante string) requestresponse.APIResponse {
	var resultado = make(map[string]interface{})
	var persona []map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{}
	var success bool = true
	var statusCode int = 200

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero?query=Id:"+idSolicitante, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		var tipoVinculacion []map[string]interface{}
		var correoElectronico []map[string]interface{}
		var correoInstitucional []map[string]interface{}
		var correoPersonal []map[string]interface{}
		var telefono []map[string]interface{}
		var celular []map[string]interface{}
		var jsondata map[string]interface{}

		// Correo institucional --> 94
		errCorreoIns := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idSolicitante)+",InfoComplementariaId__Id:94,Activo:true", &correoInstitucional)
		if errCorreoIns == nil && fmt.Sprintf("%v", correoInstitucional[0]) != "map[]" {
			if correoInstitucional[0]["Status"] != 404 {
				correoaux := correoInstitucional[0]["Dato"]
				if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
					panic(err)
				}
				resultado["CorreoInstitucional"] = jsondata["value"]
			} else {
				if correoInstitucional[0]["Message"] == "Not found resource" {

				} else {
					alertas = append(alertas, "No data found")
					alerta.Code = 404
					alerta.Type = "error"
				}
			}
		} else {
			alertas = append(alertas, "No data found")
			alerta.Code = 400
			alerta.Type = "error"
		}

		// Correo --> 53
		errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idSolicitante)+",InfoComplementariaId__Id:53,Activo:true", &correoElectronico)
		if errCorreo == nil && fmt.Sprintf("%v", correoElectronico[0]) != "map[]" {
			if correoElectronico[0]["Status"] != 404 {
				correoaux := correoElectronico[0]["Dato"]
				if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
					panic(err)
				}
				resultado["Correo"] = jsondata["Data"]
			} else {
				if correoElectronico[0]["Message"] == "Not found resource" {

				} else {
					alertas = append(alertas, "No data found")
					alerta.Code = 404
					alerta.Type = "error"
				}
			}
		} else {
			alertas = append(alertas, "No data found")
			alerta.Code = 400
			alerta.Type = "error"
		}

		// Correo personal --> 253
		errCorreoPersonal := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idSolicitante)+",InfoComplementariaId__Id:253,Activo:true", &correoPersonal)
		if errCorreoPersonal == nil && fmt.Sprintf("%v", correoPersonal[0]) != "map[]" {
			if correoPersonal[0]["Status"] != 404 {
				correoaux := correoPersonal[0]["Dato"]
				if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
					panic(err)
				}
				resultado["CorreoPersonal"] = jsondata["Data"]
			} else {
				if correoPersonal[0]["Message"] == "Not found resource" {

				} else {
					alertas = append(alertas, "No data found")
					alerta.Code = 404
					alerta.Type = "error"
				}
			}
		} else {
			alertas = append(alertas, "No data found")
			alerta.Code = 400
			alerta.Type = "error"
		}

		// Teléfono --> 51
		errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idSolicitante)+",InfoComplementariaId__Id:51,Activo:true", &telefono)
		if errTelefono == nil && fmt.Sprintf("%v", telefono[0]) != "map[]" {
			if telefono[0]["Status"] != 404 {
				telefonoaux := telefono[0]["Dato"]

				if err := json.Unmarshal([]byte(telefonoaux.(string)), &jsondata); err != nil {
					resultado["Telefono"] = telefono[0]["Dato"]
				} else {
					resultado["Telefono"] = jsondata["principal"]
				}
			} else {
				if telefono[0]["Message"] == "Not found resource" {

				} else {
					alertas = append(alertas, "No data found")
					alerta.Code = 404
					alerta.Type = "error"
				}
			}
		} else {
			alertas = append(alertas, "No data found")
			alerta.Code = 400
			alerta.Type = "error"
		}

		// Celular --> 52
		errCelular := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idSolicitante)+",InfoComplementariaId__Id:52,Activo:true", &celular)
		if errCelular == nil && fmt.Sprintf("%v", celular[0]) != "map[]" {
			if celular[0]["Status"] != 404 {
				celularaux := celular[0]["Dato"]

				if err := json.Unmarshal([]byte(celularaux.(string)), &jsondata); err != nil {
					resultado["Celular"] = celular[0]["Dato"]
				} else {
					resultado["Celular"] = jsondata["principal"]
				}
			} else {
				if celular[0]["Message"] == "Not found resource" {

				} else {
					alertas = append(alertas, "No data found")
					alerta.Code = 404
					alerta.Type = "error"
				}
			}
		} else {
			alertas = append(alertas, "No data found")
			alerta.Code = 400
			alerta.Type = "error"
		}

		// DOCENTE DE PLANTA 	292
		// DOCENTE DE CARRERA TIEMPO COMPLETO 	293
		// DOCENTE DE CARRERA MEDIO TIEMPO 	294
		// DOCENTE DE VINCULACIÓN ESPECIAL 	295
		// HORA CÁTEDRA 	297
		// TIEMPO COMPLETO OCASIONAL 	296
		// MEDIO TIEMPO OCASIONAL 	298
		// HORA CÁTEDRA POR HONORARIOS 	299
		errVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"vinculacion?query=TerceroPrincipalId:"+fmt.Sprintf("%v", idSolicitante)+",Activo:true&limit=0", &tipoVinculacion)
		if errVinculacion == nil && fmt.Sprintf("%v", tipoVinculacion[0]) != "map[]" {
			if tipoVinculacion[0]["Status"] != 404 {
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
								resultado["TipoVinculacionId"] = vinculacion["Data"].([]interface{})[0]
							}
						}
					}
				}
			} else {
				if tipoVinculacion[0]["Message"] == "Not found resource" {

				} else {
					alertas = append(alertas, "No data found")
					alerta.Code = 404
					alerta.Type = "error"
					alerta.Body = alertas
					success = false
				}
			}
		} else {
			alertas = append(alertas, "No data found")
			success = false
			alerta.Code = 400
			alerta.Body = alertas
			alerta.Type = "error"
		}

		resultado["Nombre"] = persona[0]["NombreCompleto"]
		resultado["Id"], _ = strconv.ParseInt(idSolicitante, 10, 64)
		resultado["PuedeBorrar"] = false
	} else {
		logs.Error(errPersona)
		statusCode = 400
		alerta.Body = errPersona.Error()
		success = false
	}

	if success {
		return requestresponse.APIResponseDTO(success, statusCode, resultado)
	} else {
		return requestresponse.APIResponseDTO(success, statusCode, nil, alerta.Body)
	}
}

func ConsultarInfoColaborador(idColaborador string) requestresponse.APIResponse {
	var resultado = make(map[string]interface{})
	var persona []map[string]interface{}
	var alerta models.Alert
	alertas := []interface{}{}

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Numero:"+idColaborador, &persona)
	fmt.Println("http://" + beego.AppConfig.String("TercerosService") + "datos_identificacion?query=Numero:" + idColaborador)

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
			errVinculacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"vinculacion?query=TerceroPrincipalId:"+fmt.Sprintf("%v", idTercero)+",Activo:true&limit=0", &tipoVinculacion)
			if errVinculacion == nil && fmt.Sprintf("%v", tipoVinculacion[0]) != "map[]" {
				if tipoVinculacion[0]["Status"] != 404 {

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
									resultado["TipoVinculacionId"] = vinculacion["Data"].([]interface{})[0]
								}
							}

							// Correo institucional --> 94
							errCorreoIns := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:94,Activo:true", &correoInstitucional)
							if errCorreoIns == nil && fmt.Sprintf("%v", correoInstitucional[0]) != "map[]" {
								if correoInstitucional[0]["Status"] != 404 {
									correoaux := correoInstitucional[0]["Dato"]
									if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
										panic(err)
									}
									resultado["CorreoInstitucional"] = jsondata["value"]
								} else {
									if correoInstitucional[0]["Message"] == "Not found resource" {

									} else {
										alertas = append(alertas, "No data found")
										alerta.Code = 404
										alerta.Type = "error"
									}
								}
							} else {
								alertas = append(alertas, "No data found")
								alerta.Code = 400
								alerta.Type = "error"
							}

							// Correo --> 53
							errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:53,Activo:true", &correoElectronico)
							if errCorreo == nil && fmt.Sprintf("%v", correoElectronico[0]) != "map[]" {
								if correoElectronico[0]["Status"] != 404 {
									correoaux := correoElectronico[0]["Dato"]
									if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
										panic(err)
									}
									resultado["Correo"] = jsondata["Data"]
								} else {
									if correoElectronico[0]["Message"] == "Not found resource" {

									} else {
										alertas = append(alertas, "No data found")
										alerta.Code = 404
										alerta.Type = "error"
									}
								}
							} else {
								alertas = append(alertas, "No data found")
								alerta.Code = 400
								alerta.Type = "error"
							}

							// Correo personal --> 253
							errCorreoPersonal := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:253,Activo:true", &correoPersonal)
							if errCorreoPersonal == nil && fmt.Sprintf("%v", correoPersonal[0]) != "map[]" {
								if correoPersonal[0]["Status"] != 404 {
									correoaux := correoPersonal[0]["Dato"]
									if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
										panic(err)
									}
									resultado["CorreoPersonal"] = jsondata["Data"]
								} else {
									if correoPersonal[0]["Message"] == "Not found resource" {

									} else {
										alertas = append(alertas, "No data found")
										alerta.Code = 404
										alerta.Type = "error"
									}
								}
							} else {
								alertas = append(alertas, "No data found")
								alerta.Code = 400
								alerta.Type = "error"
							}

							// Teléfono --> 51
							errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:51,Activo:true", &telefono)
							if errTelefono == nil && fmt.Sprintf("%v", telefono[0]) != "map[]" {
								if telefono[0]["Status"] != 404 {
									telefonoaux := telefono[0]["Dato"]

									if err := json.Unmarshal([]byte(telefonoaux.(string)), &jsondata); err != nil {
										resultado["Telefono"] = telefono[0]["Dato"]
									} else {
										resultado["Telefono"] = jsondata["principal"]
									}
								} else {
									if telefono[0]["Message"] == "Not found resource" {

									} else {
										alertas = append(alertas, "No data found")
										alerta.Code = 404
										alerta.Type = "error"
									}
								}
							} else {
								alertas = append(alertas, "No data found")
								alerta.Code = 400
								alerta.Type = "error"
							}

							// Celular --> 52
							errCelular := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", idTercero)+",InfoComplementariaId__Id:52,Activo:true", &celular)
							if errCelular == nil && fmt.Sprintf("%v", celular[0]) != "map[]" {
								if celular[0]["Status"] != 404 {
									celularaux := celular[0]["Dato"]

									if err := json.Unmarshal([]byte(celularaux.(string)), &jsondata); err != nil {
										resultado["Celular"] = celular[0]["Dato"]
									} else {
										resultado["Celular"] = jsondata["principal"]
									}
								} else {
									if celular[0]["Message"] == "Not found resource" {

									} else {
										return requestresponse.APIResponseDTO(false, 404, nil, "No data found")
									}
								}
							} else {
								return requestresponse.APIResponseDTO(false, 404, nil, "No data found")
							}

							resultado["Nombre"] = persona[0]["TerceroId"].(map[string]interface{})["NombreCompleto"]
							resultado["Id"] = idTercero
							resultado["PuedeBorrar"] = true
							break
						} else {
							return requestresponse.APIResponseDTO(false, 404, nil, "No es docente")
						}

					}

				} else {
					if tipoVinculacion[0]["Message"] == "Not found resource" {

					} else {
						return requestresponse.APIResponseDTO(false, 404, nil, "No data found")
					}
				}
			} else {
				return requestresponse.APIResponseDTO(false, 404, nil, "No data found")
			}

		} else {
			if persona[0]["Message"] == "Not found resource" {

			} else {
				return requestresponse.APIResponseDTO(false, 404, nil, "No data found")
			}
		}
	} else {
		if errPersona == nil {
			return requestresponse.APIResponseDTO(false, 404, nil, "Error en la consulta a terceros")
		} else {
			return requestresponse.APIResponseDTO(false, 404, nil, errPersona.Error())
		}
	}

	return requestresponse.APIResponseDTO(true, 200, resultado)

}

func ConsultarParametros() requestresponse.APIResponse {
	var getProyecto []map[string]interface{}
	var proyectos []map[string]interface{}
	var estados []interface{}
	var vehiculos map[string]interface{}
	var resultado = make(map[string]interface{})
	var periodos map[string]interface{}
	var tipoEstados map[string]interface{}
	var tipoSolicitud map[string]interface{}

	errPeriodo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"periodo?query=CodigoAbreviacion:PA,Activo:true&limit=1&sortby=Id&order=desc", &periodos)
	if errPeriodo == nil && fmt.Sprintf("%v", periodos["Data"]) != "[map[]]" {
		if periodos["Status"] != "404" {
			resultado["periodos"] = periodos["Data"]
		}
	} else {
		resultado["periodos"] = nil
		logs.Error(periodos)
		return requestresponse.APIResponseDTO(false, 404, nil, errPeriodo.Error())
	}

	errProyecto := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"proyecto_academico_institucion/?query=Activo:true,Oferta:true&limit=0", &getProyecto)
	if errProyecto == nil {
		for _, proyectoAux := range getProyecto {
			proyecto := map[string]interface{}{
				"Id":          proyectoAux["Id"],
				"Nombre":      proyectoAux["Nombre"],
				"Codigo":      proyectoAux["Codigo"],
				"CodigoSnies": proyectoAux["CodigoSnies"],
			}
			proyectos = append(proyectos, proyecto)
		}
		resultado["proyectos"] = proyectos
	} else {
		resultado["proyectos"] = nil
		logs.Error(getProyecto)
		return requestresponse.APIResponseDTO(false, 404, nil, errProyecto.Error())
	}

	errVehiculo := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+"parametro/?query=tipo_parametro_id:38&sortby=numero_orden&order=asc&limit=0", &vehiculos)
	if errVehiculo == nil && fmt.Sprintf("%v", vehiculos["Data"]) != "[map[]]" {
		if vehiculos["Status"] != "404" {
			resultado["vehiculos"] = vehiculos["Data"]
		}
	} else {
		resultado["proyectos"] = nil
		logs.Error(getProyecto)
		return requestresponse.APIResponseDTO(false, 404, nil, errVehiculo.Error())
	}

	errTipoSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"tipo_solicitud?query=CodigoAbreviacion:SoPA", &tipoSolicitud)
	if errTipoSolicitud == nil && fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0]) != "map[]" {
		var id = fmt.Sprintf("%v", tipoSolicitud["Data"].([]interface{})[0].(map[string]interface{})["Id"])

		errTipoEstados := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"estado_tipo_solicitud?query=TipoSolicitud.Id:"+id, &tipoEstados)
		if errTipoEstados == nil && fmt.Sprintf("%v", tipoEstados["Data"]) != "[map[]]" {
			if tipoEstados["Status"] != "404" {
				for _, estado := range tipoEstados["Data"].([]interface{}) {
					estados = append(estados, estado.(map[string]interface{})["EstadoId"])
					resultado["estados"] = estados
				}
			}
		} else {
			logs.Error(tipoEstados)
			return requestresponse.APIResponseDTO(false, 404, nil, errTipoEstados.Error())
		}
	}
	return requestresponse.APIResponseDTO(true, 200, resultado)
}

func ConsultarEspaciosAcademicos(idProyecto string) requestresponse.APIResponse {
	resultado := []interface{}{}
	var espaciosAcademicos map[string]interface{}

	errEspaciosAcademicos := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico?query=activo:true,docente_id:"+fmt.Sprintf("%v", idProyecto), &espaciosAcademicos)
	if errEspaciosAcademicos == nil && fmt.Sprintf("%v", espaciosAcademicos["Data"]) != "[map[]]" {
		if espaciosAcademicos["Status"] != "404" {

			for _, espacioAcademico := range espaciosAcademicos["Data"].([]interface{}) {

				resultado = append(resultado, map[string]interface{}{
					"Nombre": fmt.Sprintf("%v", espacioAcademico.(map[string]interface{})["nombre"]) + " - " + fmt.Sprintf("%v", espacioAcademico.(map[string]interface{})["grupo"]),
					"Id":     espacioAcademico.(map[string]interface{})["_id"],
				})
			}
		}
	} else {
		return requestresponse.APIResponseDTO(false, 404, nil, errEspaciosAcademicos.Error())
	}
	return requestresponse.APIResponseDTO(true, 200, resultado)
}

func EnviarInvitaciones(data []byte) requestresponse.APIResponse {
	var Solicitudes []map[string]interface{}
	var CorreoPost map[string]interface{}
	var solicitud map[string]interface{}

	if err := json.Unmarshal(data, &solicitud); err == nil {
		id_practica := solicitud["Id"]

		errSolicitud := request.GetJson("http://"+beego.AppConfig.String("SolicitudDocenteService")+"solicitante?query=SolicitudId.Id:"+id_practica.(string), &Solicitudes)
		if errSolicitud == nil {
			if Solicitudes != nil && fmt.Sprintf("%v", Solicitudes[0]) != "map[]" {
				Referencia := Solicitudes[0]["SolicitudId"].(map[string]interface{})["Referencia"].(string)
				var ReferenciaJson map[string]interface{}
				if err := json.Unmarshal([]byte(Referencia), &ReferenciaJson); err == nil {
					ReferenciaJson["Id"] = id_practica
				}

				idEstado := fmt.Sprintf("%v", Solicitudes[0]["SolicitudId"].(map[string]interface{})["EstadoTipoSolicitudId"].(map[string]interface{})["Id"].(float64))

				if idEstado == "39" {

					// TO DO: Consulta de correos electronicos de los estudiantes inscritos en el espacio académico
					correoEstudiantes := []interface{}{
						"correo1@correo.com", "correo2@correo.com",
					}
					nombreEstudiantes := []interface{}{
						"Nombre 1", "Nombre 2",
					}

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

						errEnvioCorreos := request.SendJson("http://"+beego.AppConfig.String("GOOGLE_MID")+"notificacion", "POST", &CorreoPost, correo)
						if errEnvioCorreos == nil {
							if CorreoPost == nil || fmt.Sprintf("%v", CorreoPost) == "400" {
								return requestresponse.APIResponseDTO(false, 400, nil, "Error en el envio de correos")
							}
						} else {
							return requestresponse.APIResponseDTO(false, 400, nil, errEnvioCorreos.Error())
						}

					}
				}

			} else {
				return requestresponse.APIResponseDTO(false, 404, nil, "No data found")
			}
		} else {
			return requestresponse.APIResponseDTO(false, 400, nil, errSolicitud.Error())
		}

	}
	return requestresponse.APIResponseDTO(true, 200, nil)

}
