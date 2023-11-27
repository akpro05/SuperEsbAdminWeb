const BASE_URL = 'http://localhost:5010';

export const LOGIN_URL = `${BASE_URL}/Login`;

export const DASHBOARD_URL = `${BASE_URL}/Dashboard`;

export const CHANGEPASSWORD_URL = `${BASE_URL}/ChangePassword`;
export const FORGOT_URL = `${BASE_URL}/ForgotPassword`;

export const SEARCH_USER_URL = `${BASE_URL}/SysUser/SearchsysUser`;
export const CREATE_USER_URL = `${BASE_URL}/SysUser/CreatesysUser`;
export const VIEW_USER_URL = `${BASE_URL}/SysUser/ViewsysUser/:Id`;
export const UPDATE_USER_URL = `${BASE_URL}/SysUser/UpdatesysUser/:Id`;

export const SEARCH_CONSUMERS_URL = `${BASE_URL}/Consumers/SearchConsumers`;
export const CREATE_CONSUMERS_URL = `${BASE_URL}/Consumers/CreateConsumers`;
export const VIEW_CONSUMERS_URL = `${BASE_URL}/Consumers/ViewConsumers/:Id`;
export const UPDATE_CONSUMERS_URL = `${BASE_URL}/Consumers/UpdateConsumers/:Id`;

export const SEARCH_PRODUCERS_URL = `${BASE_URL}/Producers/SearchProducers`;
export const CREATE_PRODUCERS_URL = `${BASE_URL}/Producers/CreateProducers`;
export const VIEW_PRODUCERS_URL = `${BASE_URL}/Producers/ViewProducers/:Id`;
export const UPDATE_PRODUCERS_URL = `${BASE_URL}/Producers/UpdateProducers/:Id`;

export const SEARCH_PROD2CONS_URL = `${BASE_URL}/ProducerToConsumer/SearchProducerToConsumer`;
export const CREATE_PROD2CONS_URL = `${BASE_URL}/ProducerToConsumer/CreateProducerToConsumer`;
export const VIEW_PROD2CONS_URL = `${BASE_URL}/ProducerToConsumer/ViewProducerToConsumer/:Id`;
export const UPDATE_PROD2CONS_URL = `${BASE_URL}/ProducerToConsumer/UpdateProducerToConsumer/:Id`;
export const GET_CONSUMER_SERVICES_URL = `${BASE_URL}/ProducerToConsumer/CreateProducerToConsumer/GetServiceList`;


export const PRODUCER_REPORT_URL = `${BASE_URL}/Reports/ProducerReport`;
export const CONSUMER_REPORT_URL = `${BASE_URL}/Reports/ConsumerReport`;
export const ESBLOGS_REPORT_URL = `${BASE_URL}/Reports/ESBLogsReport`;
export const AUDIT_REPORT_URL = `${BASE_URL}/Reports/AuditReport`;

export const SEARCH_ROLE_URL = `${BASE_URL}/Role/SearchRole`;
export const CREATE_ROLE_URL = `${BASE_URL}/Role/CreateRole`;
export const VIEW_ROLE_URL = `${BASE_URL}/Role/ViewRole/:Id`;
export const UPDATE_ROLE_URL = `${BASE_URL}/Role/UpdateRole/:Id`;