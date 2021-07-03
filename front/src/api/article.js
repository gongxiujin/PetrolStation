import request from '@/utils/request'

export function fetchList(query) {
  return request({
    url: '/management/users',
    method: 'get',
    params: query
  })
}

export function dashboardchart() {
  return request({
    url: '/management/dashboard_char',
    method: 'get',
  })
}

export function get_statistical() {
  return request({
    url: '/management/get_statistical',
    method: 'get',
  })
}

export function delete_user(id) {
  return request({
    url: '/management/users',
    method: 'delete',
    params: {"user_id": id}
  })
}

export function update_user_figure(id, data) {
  return request({
    url: '/management/users/'+ id,
    method: 'post',
    data:data
  })
}



export function user_profile(id) {
  return request({
    url: '/management/users/'+id,
    method: 'get',
  })
}

export function moments(query) {
  return request({
    url: '/management/discover',
    method: 'get',
    params: query
  })
}

export function delete_moment(id) {
  return request({
    url: '/management/discover',
    method: 'delete',
    data:{"id": id}
  })
}

export function get_comments(query) {
  return request({
    url: '/management/comments',
    method: 'get',
    params: query
  })
}

export function delete_comments(id) {
  return request({
    url: '/management/comments',
    method: 'delete',
    data: {"id": id}
  })
}

export function get_reports(query) {
  return request({
    url: '/management/reports',
    method: 'get',
    params:query
  })
}
export function handle_reports(id, result) {
  return request({
    url: '/management/reports',
    method: 'post',
    data: {"id": id, "result":result}
  })
}
