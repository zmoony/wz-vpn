import { createApp } from "vue";
import { createPinia } from "pinia";
import { ElAlert } from "element-plus/es/components/alert/index";
import { ElButton } from "element-plus/es/components/button/index";
import { ElDescriptions, ElDescriptionsItem } from "element-plus/es/components/descriptions/index";
import { ElDialog } from "element-plus/es/components/dialog/index";
import { ElDivider } from "element-plus/es/components/divider/index";
import { ElEmpty } from "element-plus/es/components/empty/index";
import { ElForm, ElFormItem } from "element-plus/es/components/form/index";
import { ElInput } from "element-plus/es/components/input/index";
import { ElInputNumber } from "element-plus/es/components/input-number/index";
import { ElOption } from "element-plus/es/components/select/index";
import { ElSelect } from "element-plus/es/components/select/index";
import { ElSpace } from "element-plus/es/components/space/index";
import { ElSwitch } from "element-plus/es/components/switch/index";
import { ElTable, ElTableColumn } from "element-plus/es/components/table/index";
import { ElTag } from "element-plus/es/components/tag/index";
import "element-plus/es/components/alert/style/css";
import "element-plus/es/components/button/style/css";
import "element-plus/es/components/descriptions/style/css";
import "element-plus/es/components/dialog/style/css";
import "element-plus/es/components/divider/style/css";
import "element-plus/es/components/empty/style/css";
import "element-plus/es/components/form/style/css";
import "element-plus/es/components/form-item/style/css";
import "element-plus/es/components/input/style/css";
import "element-plus/es/components/input-number/style/css";
import "element-plus/es/components/option/style/css";
import "element-plus/es/components/select/style/css";
import "element-plus/es/components/space/style/css";
import "element-plus/es/components/switch/style/css";
import "element-plus/es/components/table/style/css";
import "element-plus/es/components/tag/style/css";

import App from "./App.vue";
import router from "./router";
import "./styles.css";

const app = createApp(App);

app.use(createPinia());
app.use(router);

[
  ElAlert,
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElDialog,
  ElDivider,
  ElEmpty,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElOption,
  ElSelect,
  ElSpace,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTag,
].forEach((component) => {
  app.component(component.name!, component);
});

app.mount("#app");
