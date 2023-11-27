import React from 'react';
import { lazy } from 'react';

// project imports
import MainLayout from '../layout/MainLayout';
import Loadable from '../ui-component/Loadable';

// dashboard routing
const DashboardDefault = Loadable(lazy(() => import('../views/dashboard/Default')));

const ChangePassword = Loadable(lazy(() => import('../views/utilities/ChangePassword/ChangePassword')));

// utilities routing
const UtilsTablerIcons = Loadable(lazy(() => import('./../views/utilities/User/SearchUser')));
const CreateUser = Loadable(lazy(() => import('./../views/utilities/User/CreateUser')));
const ViewUser = Loadable(lazy(() => import('./../views/utilities/User/ViewUser')));
const UpdateUser = Loadable(lazy(() => import('./../views/utilities/User/UpdateUser')));

const ESBLOgsReport = Loadable(lazy(() => import('./../views/utilities/Reports/ESBLogsReport')));

const CreateConsumers = Loadable(lazy(() => import('./../views/utilities/Consumers/CreateConsumers')));
const SearchConsumers = Loadable(lazy(() => import('./../views/utilities/Consumers/SearchConsumers')));
const ViewConsumers = Loadable(lazy(() => import('./../views/utilities/Consumers/ViewConsumers')));
const UpdateConsumers = Loadable(lazy(() => import('./../views/utilities/Consumers/UpdateConsumers')));


const CreateProducers = Loadable(lazy(() => import('./../views/utilities/Producers/CreateProducers')));
const SearchProducers = Loadable(lazy(() => import('./../views/utilities/Producers/SearchProducers')));
const ViewProducers = Loadable(lazy(() => import('./../views/utilities/Producers/ViewProducers')));
const UpdateProducers = Loadable(lazy(() => import('./../views/utilities/Producers/UpdateProducers')));


const CreateProducerToConsumer = Loadable(lazy(() => import('./../views/utilities/SystemConfig/ProducerToConsumer/CreateProducerToConsumer')));
const SearchProducerToConsumer = Loadable(lazy(() => import('./../views/utilities/SystemConfig/ProducerToConsumer/SearchProducerToConsumer')));
const ViewProducerToConsumer = Loadable(lazy(() => import('./../views/utilities/SystemConfig/ProducerToConsumer/ViewProducerToConsumer')));
const UpdateProducerToConsumer = Loadable(lazy(() => import('./../views/utilities/SystemConfig/ProducerToConsumer/UpdateProducerToConsumer')));

// sample page routing
const SamplePage = Loadable(lazy(() => import('../views/sample-page')));

const ProducerReport = Loadable(lazy(() => import('./../views/utilities/Reports/ProducerReport')));
const ConsumerReport = Loadable(lazy(() => import('./../views/utilities/Reports/ConsumerReport')));
const AuditReport = Loadable(lazy(() => import('./../views/utilities/Reports/AuditReport')));

const CreateRole = Loadable(lazy(() => import('./../views/utilities/SystemConfig/Role/CreateRole')));
const SearchRole = Loadable(lazy(() => import('./../views/utilities/SystemConfig/Role/SearchRole')));
const ViewRole = Loadable(lazy(() => import('./../views/utilities/SystemConfig/Role/ViewRole')));
const UpdateRole = Loadable(lazy(() => import('./../views/utilities/SystemConfig/Role/UpdateRole')));



// ==============================|| MAIN ROUTING ||============================== //

const MainRoutes = {
  path: '/',
  element: <MainLayout />,
  children: [
    { path: '/Dashboard', element: <DashboardDefault /> },
    { path: '/ChangePassword', element: <ChangePassword /> },

    { path: 'SysUser', children: [{ path: 'SearchsysUser', element: <UtilsTablerIcons /> }] },
    { path: 'SysUser', children: [{ path: 'CreatesysUser', element: <CreateUser /> }] },
    { path: 'SysUser', children: [{ path: 'ViewsysUser/:id', element: <ViewUser /> }] },
    { path: 'SysUser', children: [{ path: 'UpdatesysUser/:id', element: <UpdateUser /> }] },

    { path: 'Consumers', children: [{ path: 'CreateConsumers', element: <CreateConsumers /> }] },
    { path: 'Consumers', children: [{ path: 'SearchConsumers', element: <SearchConsumers /> }] },
    { path: 'Consumers', children: [{ path: 'UpdateConsumers/:id', element: <UpdateConsumers /> }] },
    { path: 'Consumers', children: [{ path: 'ViewConsumers/:id', element: <ViewConsumers /> }] },

    { path: 'Producers', children: [{ path: 'CreateProducers', element: <CreateProducers /> }] },
    { path: 'Producers', children: [{ path: 'SearchProducers', element: <SearchProducers /> }] },
    { path: 'Producers', children: [{ path: 'UpdateProducers/:id', element: <UpdateProducers /> }] },
    { path: 'Producers', children: [{ path: 'ViewProducers/:id', element: <ViewProducers /> }] },

    { path: 'ProducerToConsumer', children: [{ path: 'CreateProducerToConsumer', element: <CreateProducerToConsumer /> }] },
    { path: 'ProducerToConsumer', children: [{ path: 'SearchProducerToConsumer', element: <SearchProducerToConsumer /> }] },
    { path: 'ProducerToConsumer', children: [{ path: 'UpdateProducerToConsumer/:id', element: <UpdateProducerToConsumer /> }] },
    { path: 'ProducerToConsumer', children: [{ path: 'ViewProducerToConsumer/:id', element: <ViewProducerToConsumer /> }] },

    { path: 'Reports', children: [{ path: 'ESBLogsReport', element: <ESBLOgsReport /> }] },
    { path: 'Reports', children: [{ path: 'ProducerReport', element: <ProducerReport /> }] },
    { path: 'Reports', children: [{ path: 'ConsumerReport', element: <ConsumerReport /> }] },
    { path: 'Reports', children: [{ path: 'AuditReport', element: <AuditReport /> }] },
    
    { path: 'Role', children: [{ path: 'CreateRole', element: <CreateRole /> }] },
    { path: 'Role', children: [{ path: 'SearchRole', element: <SearchRole /> }] },
    { path: 'Role', children: [{ path: 'UpdateRole/:id', element: <UpdateRole /> }] },
    { path: 'Role', children: [{ path: 'ViewRole/:id', element: <ViewRole /> }] },
  ]
  
};

export default MainRoutes;