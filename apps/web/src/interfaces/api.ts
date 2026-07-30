export type ApiErrorResponse = {
  error?: string;
  message?: string;
};

export type Entity = {
  id: string;
  createdAt: string;
  updatedAt: string;
};

export type PaginatedResponse<T> = {
  items: T[];
  total: number;
};

export type PageQuery = {
  limit?: number;
  offset?: number;
};
