export type SmsTemplateItemProps = {
  id?: number
  name: string
  code: string
  content: string
  sign_name?: string
  provider_template_id?: string
  enabled: boolean
  remark?: string
}

export type FormProps = {
  formInline: SmsTemplateItemProps
}
