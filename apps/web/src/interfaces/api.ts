export type ApiError = { error: string };

export type PaginatedResponse<T> = {
  items: T[];
  total: number;
};
