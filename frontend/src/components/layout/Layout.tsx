import { Box } from '@mui/material';
import type { BaseComponentProps } from '../../types';
import Header from './Header';
import Footer from './Footer';

interface LayoutProps extends BaseComponentProps {
  children: React.ReactNode;
}

const Layout = ({ children, className = '' }: LayoutProps) => {
  return (
    <Box
      className={className}
      sx={{
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      <Header />
      <Box component="main" sx={{ flex: 1 }}>
        {children}
      </Box>
      <Footer />
    </Box>
  );
};

export default Layout;
