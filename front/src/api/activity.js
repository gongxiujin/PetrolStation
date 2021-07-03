/**
 * Created by gongxiujin on 2019/7/21.
 */
import request from '@/utils/request'

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
    headers:{"Content-Type": "multipart/form-data"},
    method: 'put',
    data: data
  })
}
