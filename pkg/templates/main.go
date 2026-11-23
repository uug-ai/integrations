package templates

import (
	"os"
)

func HasTemplateOnFileSystem(path string, template string) bool {
	if _, err := os.Stat(path + "/" + template); err == nil {
		return true
	} else {
		return false
	}
}

func Readfile(path string, template string) string {
	content, err := os.ReadFile(path + "/" + template)
	if err != nil {
		return err.Error()
	}
	return string(content)
}

func GetTemplate(name string) string {
	path := "/mail/templates"
	if HasTemplateOnFileSystem(path, name+".html") {
		return Readfile(path, name+".html")
	} else {
		if name == "welcome" {
			return TEMPLATE_WELCOME
		} else if name == "forgot" {
			return TEMPLATE_FORGOT
		} else if name == "activate" {
			return TEMPLATE_ACTIVATE
		} else if name == "detection" {
			return TEMPLATE_DETECTION
		} else if name == "device" {
			return TEMPLATE_DEVICE
		} else if name == "device_group" {
			return TEMPLATE_DEVICE_GROUP
		} else if name == "recording" {
			return TEMPLATE_RECORDING
		} else if name == "invitation" {
			return TEMPLATE_INVITATION
		} else if name == "disable" {
			return TEMPLATE_DISABLE
		} else if name == "highupload" {
			return TEMPLATE_HIGHUPLOAD
		} else if name == "newip" {
			return TEMPLATE_NEWIP
		} else if name == "share" {
			return TEMPLATE_SHARE
		} else if name == "share_case" {
			return TEMPLATE_SHARE_CASE
		} else if name == "share_case_otp" {
			return TEMPLATE_SHARE_CASE_OTP
		} else if name == "assign_task" {
			return TEMPLATE_TASK_ASSIGNMENT
		}
	}
	return ""
}

func GetTextTemplate(name string) string {
	path := "/mail/templates"
	if HasTemplateOnFileSystem(path, name+".txt") {
		return Readfile(path, name+".txt")
	} else {
		if name == "welcome" {
			return TEMPLATE_WELCOME_TEXT
		} else if name == "forgot" {
			return TEMPLATE_FORGOT_TEXT
		} else if name == "activate" {
			return TEMPLATE_ACTIVATE_TEXT
		} else if name == "detection" {
			return TEMPLATE_DETECTION_TEXT
		} else if name == "device" {
			return TEMPLATE_DEVICE_TEXT
		} else if name == "device_group" {
			return TEMPLATE_DEVICE_GROUP_TEXT
		} else if name == "recording" {
			return TEMPLATE_RECORDING_TEXT
		} else if name == "invitation" {
			return TEMPLATE_INVITATION_TEXT
		} else if name == "disable" {
			return TEMPLATE_DISABLE_TEXT
		} else if name == "highupload" {
			return TEMPLATE_HIGHUPLOAD_TEXT
		} else if name == "newip" {
			return TEMPLATE_NEWIP_TEXT
		} else if name == "share" {
			return TEMPLATE_SHARE_TEXT
		} else if name == "share_case" {
			return TEMPLATE_SHARE_CASE_TEXT
		} else if name == "share_case_otp" {
			return TEMPLATE_SHARE_CASE_OTP_TEXT
		} else if name == "assign_task" {
			return TEMPLATE_TASK_ASSIGNMENT_TEXT
		}
	}
	return ""
}
