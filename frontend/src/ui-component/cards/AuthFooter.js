import React from 'react';
// material-ui
import { Link, Typography, Stack } from '@mui/material';

// ==============================|| FOOTER - AUTHENTICATION 2 & 3 ||============================== //

const AuthFooter = () => (
  <Stack direction="row" justifyContent="space-between">
    <Typography variant="subtitle2" component={Link} href="#" target="_blank" underline="hover" style={{ color: 'white' }}>
      Supernet Technologies
    </Typography>

    <Typography variant="subtitle2" component={Link} href="https://supernet-india.com/" target="_blank" underline="hover" style={{ color: 'white' }}>
      &copy; supernet-india.com
    </Typography>
  </Stack>
);

export default AuthFooter;
