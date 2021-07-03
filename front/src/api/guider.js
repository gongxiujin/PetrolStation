/**
 * Created by gongxiujin on 2019/7/4.
 */
import request from '@/utils/request'

export function getGuider(query) {
  return request({
    url: '/management/guider',
    method: 'get',
    params: query
  })
}

export function beGuider(ID) {
  return request({
    url: '/management/guider',
    method: 'post',
    data:{"ID": ID}
  })
}

export function getRewards() {
  return request({
    url: '/management/rewards',
    method: 'get',
  })
}

export function handleReward(id) {
  return request({
    url: '/management/rewards',
    method: 'put',
    params: {"id": id},
  })
}



export function deleteReward(id) {
  return request({
    url: "/management/rewards",
    method: 'delete',
    params: {"id": id},
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
