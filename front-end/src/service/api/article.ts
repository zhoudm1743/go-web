import { http } from '../http';

// 文章列表查询参数
export interface ArticleQueryParams {
  page?: number;
  pageSize?: number;
  title?: string;
  categoryID?: string;
  author?: string;

}

// 文章创建请求
export interface ArticleCreateRequest {
  title: string;
  categoryID?: string;
  author?: string;

}

// 文章更新请求
export interface ArticleUpdateRequest {
  id: number;
  title?: string;
  categoryID?: string;
  author?: string;

}

// 文章响应
export interface ArticleResponse {
  id: number;
  createdAt: string;
  updatedAt: string;
  title: string;
  categoryID: string;
  author: string;

}

// 文章列表响应
export interface ArticleListResponse {
  total: number;
  list: ArticleResponse[];
}

// 获取文章列表
export const getArticles = (params: ArticleQueryParams) => {
  return http.request<ArticleListResponse>({
    url: '/article/list',
    method: 'GET',
    params,
  });
};

// 获取文章详情
export const getArticle = (id: number) => {
  return http.request<ArticleResponse>({
    url: '/article/detail/' + id,
    method: 'GET',
  });
};

// 创建文章
export const createArticle = (data: ArticleCreateRequest) => {
  return http.request<void>({
    url: '/article/create',
    method: 'POST',
    data,
  });
};

// 更新文章
export const updateArticle = (data: ArticleUpdateRequest) => {
  return http.request<void>({
    url: '/article/update',
    method: 'PUT',
    data,
  });
};

// 删除文章
export const deleteArticle = (id: number) => {
  return http.request<void>({
    url: '/article/delete/' + id,
    method: 'DELETE',
  });
};
