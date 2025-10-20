import {
  Box,
  Container,
  Grid,
  Typography,
  Link,
  Divider,
  useTheme
} from '@mui/material';
import { useLanguage } from '../../contexts/LanguageContext';
import type { BaseComponentProps } from '../../types';

type FooterProps = BaseComponentProps;

const Footer = ({ className = '' }: FooterProps) => {
  const theme = useTheme();
  const { t } = useLanguage();
  const currentYear = new Date().getFullYear();

  return (
    <Box
      component="footer"
      sx={{
        backgroundColor: theme.palette.mode === 'light'
          ? 'rgba(255, 255, 255, 0.9)'
          : 'rgba(30, 30, 30, 0.9)',
        backdropFilter: 'blur(10px)',
        borderTop: `1px solid ${theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.15)' : theme.palette.divider}`,
        mt: 'auto',
      }}
      className={className}
    >
      <Container maxWidth="lg" sx={{ py: 6 }}>
        <Grid container spacing={6} justifyContent="center">
          <Grid item xs={12} md={4} sx={{ textAlign: { xs: 'center', md: 'left' } }}>
            <Typography
              variant="h5"
              sx={{
                fontWeight: 700,
                color: theme.palette.primary.main,
                mb: 2,
              }}
            >
              MCD
            </Typography>
            <Typography
              variant="body1"
              color="text.secondary"
              sx={{
                maxWidth: 300,
                mx: { xs: 'auto', md: 0 },
                lineHeight: 1.6
              }}
            >
              {t('footer.description')}
            </Typography>
          </Grid>

          <Grid item xs={12} md={4} sx={{ textAlign: { xs: 'center', md: 'left' } }}>
            <Typography
              variant="h6"
              sx={{
                fontWeight: 600,
                mb: 3,
                color: theme.palette.text.primary
              }}
            >
              {t('footer.navigation')}
            </Typography>
            <Box sx={{
              display: 'flex',
              flexDirection: 'column',
              gap: 2,
              alignItems: { xs: 'center', md: 'flex-start' }
            }}>
              <Link
                href="/"
                color="text.secondary"
                underline="hover"
                sx={{
                  fontSize: '0.95rem',
                  fontWeight: 500,
                  transition: 'color 0.2s ease',
                  '&:hover': {
                    color: theme.palette.primary.main
                  }
                }}
              >
                {t('footer.home')}
              </Link>
              <Link
                href="/documents"
                color="text.secondary"
                underline="hover"
                sx={{
                  fontSize: '0.95rem',
                  fontWeight: 500,
                  transition: 'color 0.2s ease',
                  '&:hover': {
                    color: theme.palette.primary.main
                  }
                }}
              >
                {t('footer.documents')}
              </Link>
              <Link
                href="/groups"
                color="text.secondary"
                underline="hover"
                sx={{
                  fontSize: '0.95rem',
                  fontWeight: 500,
                  transition: 'color 0.2s ease',
                  '&:hover': {
                    color: theme.palette.primary.main
                  }
                }}
              >
                {t('footer.groups')}
              </Link>
              <Link
                href="/settings"
                color="text.secondary"
                underline="hover"
                sx={{
                  fontSize: '0.95rem',
                  fontWeight: 500,
                  transition: 'color 0.2s ease',
                  '&:hover': {
                    color: theme.palette.primary.main
                  }
                }}
              >
                {t('footer.settings')}
              </Link>
            </Box>
          </Grid>
        </Grid>
      </Container>

      <Divider sx={{ my: 2 }} />

      <Container maxWidth="lg" sx={{ py: 3 }}>
        <Typography
          variant="body2"
          color="text.secondary"
          align="center"
          sx={{
            fontSize: '0.9rem',
            fontWeight: 400,
            opacity: 0.8
          }}
        >
          &copy; {currentYear} MCD. {t('footer.copyright')}
        </Typography>
      </Container>
    </Box>
  );
};

export default Footer;
