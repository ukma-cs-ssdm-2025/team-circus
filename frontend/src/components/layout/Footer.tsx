import { 
  Box, 
  Container, 
  Grid, 
  Typography, 
  Link, 
  Divider,
  useTheme
} from '@mui/material';
import type { BaseComponentProps } from '../../types';

interface FooterProps extends BaseComponentProps {}

const Footer = ({ className = '' }: FooterProps) => {
  const theme = useTheme();
  const currentYear = new Date().getFullYear();

  return (
    <Box
      component="footer"
      sx={{
        backgroundColor: 'rgba(255, 255, 255, 0.9)',
        backdropFilter: 'blur(10px)',
        borderTop: `1px solid ${theme.palette.divider}`,
        mt: 'auto',
      }}
      className={className}
    >
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Grid container spacing={4}>
          <Grid item xs={12} sm={6} md={4}>
            <Typography
              variant="h6"
              sx={{
                fontWeight: 700,
                color: theme.palette.primary.main,
                mb: 2,
              }}
            >
              MCD
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Ваш надійний помічник для управління документами
            </Typography>
          </Grid>
          
          <Grid item xs={12} sm={6} md={4}>
            <Typography variant="h6" gutterBottom>
              Навігація
            </Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
              <Link href="/" color="text.secondary" underline="hover">
                Головна
              </Link>
              <Link href="/documents" color="text.secondary" underline="hover">
                Документи
              </Link>
              <Link href="/groups" color="text.secondary" underline="hover">
                Групи
              </Link>
              <Link href="/settings" color="text.secondary" underline="hover">
                Налаштування
              </Link>
            </Box>
          </Grid>
          
          <Grid item xs={12} sm={6} md={4}>
            <Typography variant="h6" gutterBottom>
              Підтримка
            </Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
              <Link href="/help" color="text.secondary" underline="hover">
                Допомога
              </Link>
              <Link href="/contact" color="text.secondary" underline="hover">
                Контакти
              </Link>
              <Link href="/privacy" color="text.secondary" underline="hover">
                Конфіденційність
              </Link>
            </Box>
          </Grid>
        </Grid>
      </Container>
      
      <Divider />
      
      <Container maxWidth="lg" sx={{ py: 2 }}>
        <Typography variant="body2" color="text.secondary" align="center">
          &copy; {currentYear} MCD. Всі права захищені.
        </Typography>
      </Container>
    </Box>
  );
};

export default Footer;
