package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "ConsultarInfoColaborador",
            Router: "/colaboradores/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "ConsultarEspaciosAcademicos",
            Router: "/espacios-academicos/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "EnviarInvitaciones",
            Router: "/invitacion/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "ConsultarParametros",
            Router: "/parametros/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_practica_academica_mid/controllers:PracticasAcademicasController"],
        beego.ControllerComments{
            Method: "ConsultarInfoSolicitante",
            Router: "/solicitantes/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
