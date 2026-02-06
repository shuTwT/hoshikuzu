import type { ExtendProForm,  ProSearchFormColumns } from "pro-naive-ui";

export type CreateProSearchFormReturn<Values = any> = ExtendProForm<Values, {
    /**
     * 是否收起
     */
    collapsed: Ref<boolean>;
    /**
     * 切换收起
     * @param collapsed 传递了此参数，根据参数切换
     */
    toggleCollapse: (collapsed?: boolean) => void;
}>;

export type CrudLayoutProps<Values = any> = {
  loading: boolean;
  form: CreateProSearchFormReturn<Values>;
  columns: ProSearchFormColumns<Values>;
}
