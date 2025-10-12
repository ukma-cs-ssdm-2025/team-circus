import { Button, useTheme } from '@mui/material';
import type { BaseComponentProps } from '../../types';

interface ActionButtonProps extends BaseComponentProps {
  children: React.ReactNode;
  onClick?: () => void;
  variant?: 'contained' | 'outlined' | 'text';
  size?: 'small' | 'medium' | 'large';
  startIcon?: React.ReactNode;
  endIcon?: React.ReactNode;
  fullWidth?: boolean;
  disabled?: boolean;
  component?: React.ElementType;
  to?: string;
}

const ActionButton = ({ 
  children, 
  onClick,
  variant = 'contained',
  size = 'large',
  startIcon,
  endIcon,
  fullWidth = false,
  disabled = false,
  component,
  to,
  className = ''
}: ActionButtonProps) => {
  const theme = useTheme();

  const getButtonStyles = () => {
    if (variant === 'contained') {
      return {
        backgroundColor: theme.palette.primary.main,
        '&:hover': {
          backgroundColor: theme.palette.primary.dark,
        }
      };
    }
    return {};
  };

  return (
    <Button
      variant={variant}
      size={size}
      onClick={onClick}
      startIcon={startIcon}
      endIcon={endIcon}
      fullWidth={fullWidth}
      disabled={disabled}
      component={component}
      to={to}
      className={className}
      sx={{
        minWidth: fullWidth ? 'auto' : 200,
        height: size === 'large' ? 60 : size === 'medium' ? 48 : 36,
        fontSize: size === 'large' ? '1.1rem' : size === 'medium' ? '1rem' : '0.875rem',
        fontWeight: 600,
        textTransform: 'none',
        borderRadius: 2,
        ...getButtonStyles(),
      }}
    >
      {children}
    </Button>
  );
};

export default ActionButton;
