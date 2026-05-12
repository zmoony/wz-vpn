import { defineStore } from "pinia";

import { fetchMe, login, logout, type SessionUser } from "../api/auth";

type AuthState = {
  initialized: boolean;
  user: SessionUser | null;
};

export const useAuthStore = defineStore("auth", {
  state: (): AuthState => ({
    initialized: false,
    user: null,
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.user),
  },
  actions: {
    async signIn(username: string, password: string) {
      const response = await login({ username, password });
      this.user = response.user;
      this.initialized = true;
    },
    async restore() {
      try {
        this.user = await fetchMe();
      } catch {
        this.user = null;
      } finally {
        this.initialized = true;
      }
    },
    async signOut() {
      await logout();
      this.user = null;
      this.initialized = true;
    },
  },
});
