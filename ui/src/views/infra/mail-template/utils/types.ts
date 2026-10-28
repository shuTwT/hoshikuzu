export type MailTemplateItemProps = {
  id?: number
  name: string
  code: string
  subject: string
  content: string
  content_type: string
  enabled: boolean
  remark?: string
}

export type FormProps = {
  formInline: MailTemplateItemProps
}
