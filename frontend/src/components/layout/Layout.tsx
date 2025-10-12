import { Box, useTheme } from '@mui/material';
import type { BaseComponentProps } from '../../types';
import Header from './Header';
import Footer from './Footer';

interface LayoutProps extends BaseComponentProps {
  children: React.ReactNode;
}

const Layout = ({ children, className = '' }: LayoutProps) => {
  const theme = useTheme();
  
  return (
    <Box
      className={className}
      sx={{
        minHeight: '100vh',
        background: theme.palette.mode === 'light' 
          ? 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)'
          : 'linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%)',
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
