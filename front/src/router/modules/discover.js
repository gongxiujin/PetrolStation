/**
 * Created by gongxiujin on 2019/7/2.
 */

/** When your routing table is too long, you can split it into small modules **/

import Layout from '@/layout'

const discoverRouter = {
  path: '/moments',
  component: Layout,
  redirect: '/moments/complex-table',
  name: '动态',
  meta: {
    title: '动态',
    icon: 'message'
  },
  children: [
    {
      path: 'list',
      component: () => import('@/views/discover/complex-table'),
      name: '动态列表',
        meta: { title: '动态列表', icon: 'guide'},
      },
    {
      path: 'details/:id(.*)',
      component: () => import('@/views/discover/details'),
      name: '动态详情',
      meta: { title: '动态详情', noCache:true, activeMenu:'/discover/list'},
      hidden: true
    },
    {
      path: 'comments',
      component: () => import('@/views/discover/components/Comments'),
      name: '评论列表', meta: { title: '评论列表',icon: 'comment'}
    },
    {
      path: 'reports',
      component: () => import('@/views/discover/components/reports'),
      name: '举报管理',
      meta: { title: '举报管理', icon: 'report'}
    }
]
}
export default discoverRouter
