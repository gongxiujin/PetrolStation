/**
 * Created by gongxiujin on 2019/7/4.
 */
import request from '@/utils/request'

export function getFace(query) {
  return request({
    url: '/management/faces',
    method: 'get',
    params: query
  })
}

export function faceSelect() {
  return request({
    url: '/management/faces',
    method: 'options',
  })
}

export function deleteImage(image_id) {
  return request({
    url: '/management/faces',
    method: 'delete',
    data:{"id": image_id}
  })
}

export function getAudit(query) {
  return request({
    url: '/management/audit',
    method: 'get',
    params: query
  })
}

export function handleVideo(id, status) {
  return request({
    url: '/management/audit',
    method: 'options',
    params: {"user_id": id},
    data: {"operation": status}
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
