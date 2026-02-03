import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Layout } from './components';
import { HomePage, PoolsPage, PoolDetailsPage, OpportunitiesPage, OpportunityDetailsPage } from './pages';
import { QUERY_STALE_TIME, QUERY_CACHE_TIME, QUERY_RETRY_COUNT } from './utils/constants';

// Create a client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: QUERY_STALE_TIME,
      gcTime: QUERY_CACHE_TIME,
      retry: QUERY_RETRY_COUNT,
      refetchOnWindowFocus: false,
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Layout>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/pools" element={<PoolsPage />} />
            <Route path="/pools/:id" element={<PoolDetailsPage />} />
            <Route path="/opportunities" element={<OpportunitiesPage />} />
            <Route path="/opportunities/:id" element={<OpportunityDetailsPage />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </Layout>
      </BrowserRouter>
    </QueryClientProvider>
  );
}

function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center py-20">
      <h1 className="text-6xl font-bold text-gray-600">404</h1>
      <p className="text-xl text-gray-400 mt-4">Page not found</p>
      <a
        href="/"
        className="mt-6 px-6 py-2 bg-primary-600 hover:bg-primary-500 text-white rounded-lg transition-colors"
      >
        Go Home
      </a>
    </div>
  );
}

export default App;
