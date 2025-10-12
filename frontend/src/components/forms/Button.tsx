import { Button as MuiButton } from '@mui/material';
import type { BaseComponentProps } from '../../types';

interface ButtonProps extends BaseComponentProps {
  children: React.ReactNode;
  variant?: 'primary' | 'secondary' | 'danger';
  size?: 'small' | 'medium' | 'large';
  onClick?: () => void;
  disabled?: boolean;
  fullWidth?: boolean;
  startIcon?: React.ReactNode;
  endIcon?: React.ReactNode;
  type?: 'button' | 'submit' | 'reset';
  className?: string;
}

const Button = ({ 
  children, 
  className = '', 
  variant = 'primary', 
  size = 'medium',
  onClick,
  disabled = false,
  fullWidth = false,
  startIcon,
  endIcon,
  type = 'button'
}: ButtonProps) => {
  const getMuiVariant = () => {
    switch (variant) {
      case 'primary':
        return 'contained';
      case 'secondary':
        return 'outlined';
      case 'danger':
        return 'contained';
      default:
        return 'contained';
    }
  };

  const getMuiSize = () => {
    switch (size) {
      case 'small':
        return 'small';
      case 'medium':
        return 'medium';
      case 'large':
        return 'large';
      default:
        return 'medium';
    }
  };

  return (
    <MuiButton
      variant={getMuiVariant()}
      size={getMuiSize()}
      className={className}
      color={variant === 'danger' ? 'error' : 'primary'}
      onClick={onClick}
      disabled={disabled}
      fullWidth={fullWidth}
      startIcon={startIcon}
      endIcon={endIcon}
      type={type}
    >
      {children}
    </MuiButton>
  );
};

export default Button;
