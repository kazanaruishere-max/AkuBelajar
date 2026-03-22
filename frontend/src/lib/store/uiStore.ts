/**
 * AkuBelajar — UI Store (Zustand)
 *
 * Manages UI state:
 * - Sidebar open/collapsed
 * - Dark mode toggle
 * - Mobile detection
 */

import { create } from 'zustand';

interface UIState {
  // State
  sidebarOpen: boolean;
  darkMode: boolean;
  isMobile: boolean;

  // Actions
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  toggleDarkMode: () => void;
  setDarkMode: (dark: boolean) => void;
  setIsMobile: (mobile: boolean) => void;
}

export const useUIStore = create<UIState>((set) => ({
  sidebarOpen: true,
  darkMode: false,
  isMobile: false,

  toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
  setSidebarOpen: (open) => set({ sidebarOpen: open }),

  toggleDarkMode: () =>
    set((state) => {
      const newDark = !state.darkMode;
      if (typeof document !== 'undefined') {
        document.documentElement.setAttribute('data-theme', newDark ? 'dark' : 'light');
        localStorage.setItem('theme', newDark ? 'dark' : 'light');
      }
      return { darkMode: newDark };
    }),

  setDarkMode: (dark) => {
    if (typeof document !== 'undefined') {
      document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light');
      localStorage.setItem('theme', dark ? 'dark' : 'light');
    }
    set({ darkMode: dark });
  },

  setIsMobile: (mobile) => set({ isMobile: mobile, sidebarOpen: !mobile }),
}));
