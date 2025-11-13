import {
  Box,
  Container,
  Typography,
  Link as MuiLink,
  Divider,
  useTheme
} from '@mui/material';
import Grid from '@mui/material/Grid';
import { Link as RouterLink } from 'react-router-dom';
import { useLanguage } from '../../contexts/LanguageContext';
import type { BaseComponentProps } from '../../types';
import { ROUTES } from '../../constants';

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
          <Grid size={{ xs: 12, md: 4 }} sx={{ textAlign: { xs: 'center', md: 'left' } }}>
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

          <Grid size={{ xs: 12, md: 4 }} sx={{ textAlign: { xs: 'center', md: 'left' } }}>
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
              <MuiLink
                component={RouterLink}
                to={ROUTES.HOME}
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
              </MuiLink>
              <MuiLink
                component={RouterLink}
                to={ROUTES.DOCUMENTS}
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
              </MuiLink>
              <MuiLink
                component={RouterLink}
                to={ROUTES.GROUPS}
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
              </MuiLink>
              <MuiLink
                component={RouterLink}
                to={ROUTES.SETTINGS}
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
              </MuiLink>
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