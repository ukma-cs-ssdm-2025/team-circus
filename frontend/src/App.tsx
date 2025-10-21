import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider } from './contexts/ThemeContext';
import { LanguageProvider } from './contexts/LanguageContext';
import { Layout } from './components/layout';
import { Home, Settings, NotFound, Documents, Groups, DocumentEditor } from './pages';
import { ROUTES } from './constants';

function App() {
  return (
    <ThemeProvider>
      <LanguageProvider>
        <Router>
          <Layout>
            <Routes>
              <Route path={ROUTES.HOME} element={<Home />} />
              <Route path={ROUTES.DOCUMENTS} element={<Documents />} />
              <Route path={`${ROUTES.DOCUMENTS}/:uuid`} element={<DocumentEditor />} />
              <Route path={ROUTES.GROUPS} element={<Groups />} />
              <Route path={ROUTES.SETTINGS} element={<Settings />} />
              <Route path="*" element={<NotFound />} />
            </Routes>
          </Layout>
        </Router>
      </LanguageProvider>
    </ThemeProvider>
  );
}

export default App;
