export type PromptImageItem = {
  image_url: string
  model_name: string
  sort: number
}

export type PromptFormItemProps = {
  id?: number
  title: string
  content: string
  negative_prompt: string
  description: string
  sort: number
  images: PromptImageItem[]
  tags: string[]
}

export type PromptFormProps = {
  formInline: PromptFormItemProps
}
