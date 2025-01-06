export type UserDetail = {
  /**
   * The account ID of the user that is unique for every account.
   */
  account_id: string;

  /**
   * The authentication token for the user.
   */
  user_token: string;
}
export type AuthResponse = {
  /**
   * Response message
   */
  message: string;

  /**
   * HTTP status code
   */
  status_code: number;

  /**
   * Success indicator
   */
  success: boolean;

  /**
   * The encrypted key
   */
  encrypted_key: string;
};
