import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider } from './contexts/ThemeContext';
import { LanguageProvider } from './contexts/LanguageContext';
import { AuthProvider } from './contexts/AuthContext';
import { Layout } from './components/layout';
import { ProtectedRoute } from './components/ProtectedRoute';
import { Home, Settings, NotFound, Documents, Groups, DocumentEditor, Login, Register, AuthTest } from './pages';
import { ROUTES } from './constants';

function App() {
  return (
    <ThemeProvider>
      <LanguageProvider>
        <AuthProvider>
          <Router>
            <Routes>
              {/* Public routes */}
              <Route path={ROUTES.LOGIN} element={<Login />} />
              <Route path={ROUTES.REGISTER} element={<Register />} />
              
              {/* Protected routes */}
              <Route path="/" element={
                <ProtectedRoute>
                  <Layout>
                    <Routes>
                      <Route path={ROUTES.HOME} element={<Home />} />
                      <Route path={ROUTES.DOCUMENTS} element={<Documents />} />
                      <Route path={`${ROUTES.DOCUMENTS}/:uuid`} element={<DocumentEditor />} />
                      <Route path={ROUTES.GROUPS} element={<Groups />} />
                      <Route path={ROUTES.SETTINGS} element={<Settings />} />
                      <Route path="/auth-test" element={<AuthTest />} />
                      <Route path="*" element={<NotFound />} />
                    </Routes>
                  </Layout>
                </ProtectedRoute>
              } />
            </Routes>
          </Router>
        </AuthProvider>
      </LanguageProvider>
    </ThemeProvider>
  );
}

export default App;
