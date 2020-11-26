package NoticeBuilder

import (
	"bytes"
	"html/template"
	"vcs.torbor.ru/notizer/workers/spectacler/internal/config"
	"vcs.torbor.ru/notizer/workers/spectacler/internal/models"
	"vcs.torbor.ru/notizer/workers/spectacler/internal/storage"
)

type Builder struct {
	ts  *storage.TemplateStorage
	cfg *config.Config
}

func New(s *storage.TemplateStorage, cfg *config.Config) *Builder {
	return &Builder{ts: s, cfg: cfg}
}

func (nb Builder) ParseNotice(nt *models.Notification) (error, *bytes.Buffer) {
	err := nt.Fill()
	if err != nil {
		return err, nil
	}

	if nt.Type == "sms" {
		return nb.buildSms(nt)
	}

	return nb.buildEmail(nt)
}

func (nb Builder) buildEmail(nt *models.Notification) (error, *bytes.Buffer) {
	tmplName := nb.ts.GenerateTemplateName(nt.Template, nt.Project, nt.Type)
	tmpl, err := nb.ts.GetTemplateByName(tmplName)

	if err != nil {
		return err, nil
	}

	layout := `{{ define "layout" }}{{template "content" . }}{{ end }}`

	if len(nt.Layout) >= 1 && !nt.NoLayout {
		layout, err = nb.ts.GetTemplateByName(nb.ts.GenerateTemplateName(nt.Layout, nt.Project, nt.Type))
	}

	htmlTemplate, err := template.New("HtmlTemplate").Funcs(template.FuncMap{
		"noescape": func(str string) template.HTML {
			return template.HTML(str)
		},
	}).Parse(layout + tmpl)

	if err != nil {
		return err, nil
	}

	nt.Params["publicUrl"] = nb.ts.Swift.StorageUrl + "/" + nb.cfg.SwiftPubContainer

	var buf bytes.Buffer

	attTemplate, err := htmlTemplate.Parse(models.AttachmentTemplate)

	if err != nil {
		return err, nil
	}

	// add attachments if email has them and if it enable in config
	if nb.cfg.AppSendAttachments {
		attch := nb.ts.GetAttachments(nt.Attachments)

		if len(attch) > 0 {
			err = attTemplate.Execute(&buf, attch)

			if err != nil {
				return err, nil
			}
		}
	}

	err = htmlTemplate.ExecuteTemplate(&buf, "layout", nt.Params)

	return err, &buf
}

func (nb Builder) buildSms(nt *models.Notification) (error, *bytes.Buffer) {
	tmplName := nb.ts.GenerateTemplateName(nt.Template, nt.Project, nt.Type)
	tmpl, err := nb.ts.GetTemplateByName(tmplName)

	if err != nil {
		return err, nil
	}

	text, err := template.New("SmsTemplate").Parse(tmpl)

	if err != nil {
		return err, nil
	}

	var buf bytes.Buffer
	err = text.Execute(&buf, nt.Params)

	return err, &buf
}
