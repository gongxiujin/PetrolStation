import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/* Router Modules */
import chartsRouter from './modules/charts'
// import tableRouter from './modules/table'
// import nestedRouter from './modules/nested'
import discoverRouter from './modules/discover'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar
    noCache: true                if set true, the page will no be cached(default is false)
    affix: true                  if set true, the tag will affix in the tags-view
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/redirect',
    mode: 'history',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '/redirect/:path*',
        component: () => import('@/views/redirect/index')
      }
      ]
},
{
  path: '/login',
  component:() =>import('@/views/login/index'),
  hidden:true
},
{
  path: '/auth-redirect',
  component:() =>import('@/views/login/auth-redirect'),
  hidden:true
},
{
  path: '/404',
  component:() =>import('@/views/error-page/404'),
  hidden:true
},
{
  path: '/401',
  component: () => import('@/views/error-page/401'),
  hidden: true
},
{
  path: '/',
  component:Layout,
  redirect:'/dashboard',
  children:
  [{
      path: 'dashboard',
      component: () => import('@/views/dashboard/index'),
      name: '首页',
      meta:{title: '首页', icon:'dashboard', affix:true}
    }
  ]
},
{
  path: '/users',
  component: Layout,
  children: [
    {
      path: 'list',
      component: () => import('@/views/users/complex-table'),
      name: '用户',
      meta:{title: '用户', icon: 'peoples', affix: true}
    },
    {
      path: ':id(.*)',
        component:() =>import('@/views/users/profile'),
      name:'userprofile',
      meta:{
      title: '用户详情',
        noCache:true,
        activeMenu:'/user/list'
      },
      hidden: true
    }
  ]
},
discoverRouter,
{
  path: '/face',
  component:Layout,
  children:[
    {
      path: 'index',
      component: () => import('@/views/face/index'),
      name: 'face',
      meta:{title: '颜值管理',
      icon:'face',
      affix:true}
    }
  ]
},
{
  path: '/audits',
    component:Layout,
  children:[
  {
    path: 'index',
    component: () => import('@/views/audit/index'),
    name: 'audit',
    meta:{title: '审核管理', icon:'audit', affix:true}
  }
]
},
{
  path: '/guider',
    component: Layout,
  redirect: '/guider/index',
  name: '推荐管理',
  meta: {
  title: '推荐管理',
    icon: 'guider'
  },
  children: [
    {
      path: 'index',
      component: () => import('@/views/guider/index'),
      name: 'guider',
    meta:{title: '推荐管理', icon:'guider',  affix:true}
    },
  {
    path: 'rewards',
      component: () => import('@/views/guider/rewards'),
    name: '提现申请',
    meta: { title: '提现申请',  icon:'rewards',  affix:true},
  },
]
},
{
  path: '/activity',
    component:Layout,
  redirect:'/activity/index',
  name: '热门活动',
  meta: {
  title: '热门活动',
    icon: 'activity'
  },
  children:[
  {
    path: 'index',
    component: () => import('@/views/activity/index'),
  name: 'Activity',
    meta:{title: '热门活动', icon:'activity', noCache:true}
},
  {
    path: 'create',
      component:() =>import('@/views/activity/create'),
    name:'CreateActivity',
    meta:{
    title: '创建活动',
      noCache:true,
      activeMenu:'/activity/index'
  },
    hidden: true
  }
]
},
{
  path: '/profile',
  component:Layout,
  redirect:'/profile/index',
  hidden:true,
  children:[
    {
      path: 'index',
      component: () => import('@/views/profile/index'),
      name: 'Profile',
      meta:{title: 'Profile', icon:'user', noCache:true}
    }
]
}
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const asyncRoutes = [

/** when your routing map is too long, you can split it into small modules **/

// 404 page must be placed at the end !!!
{
  path: '*',
  redirect:'/404',
  hidden:true
}
]

const createRouter = () => new Router({
  mode: 'history', // require service support
  scrollBehavior: () => ({y: 0}),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
