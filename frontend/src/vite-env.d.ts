/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE: string;
  readonly VITE_USE_MOCK_DATA: string;
  readonly VITE_WS_BASE: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
