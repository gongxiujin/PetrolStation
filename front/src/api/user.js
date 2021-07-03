import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/management/login',
    method: 'post',
    data
  })
}

export function getInfo(token) {
  return request({
    url: '/management/user/profile',
    method: 'get',
  })
}

export function logout() {
  return request({
    url: '/user/logout',
    method: 'post'
  })
}
export function getActivity(query) {
  return request({
    url: '/management/activity',
    method: 'get',
    params: query
  })
}

export function deleteActivity(id) {
  return request({
    url: '/management/activity',
    method: 'delete',
    data:{"id": id}
  })
}

export function createActivity(data) {
  return request({
    url: '/management/activity',
    method: 'put',
    data: data
  })
}
