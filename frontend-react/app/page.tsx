'use client';

import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import FormControlLabel from '@mui/material/FormControlLabel';
import Checkbox from '@mui/material/Checkbox';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';
import { AlertColor } from '@mui/material/Alert';

// 定义 Severity 类型
type Severity = AlertColor;




function Copyright(props: any) {
  return (
    <Typography variant="body2" color="text.secondary" align="center" {...props}>
      {'Copyright © could compute group6'}
      {new Date().getFullYear()}
      {'.'}
    </Typography>
  );
}

// TODO remove, this demo shouldn't need to reset the theme.
const defaultTheme = createTheme();

export default function SignIn() {
  const [open, setOpen] = React.useState(false);
  const [message, setMessage] = React.useState('');
  const [severity, setSeverity] = React.useState<Severity>('success');


  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const data = new FormData(event.currentTarget);
    const url = 'http://localhost:8000/auth/login';
    
    try {
      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          username: data.get('username'),
          password: data.get('password'),
        })
      });
      console.log({
        email: data.get('username'),
        password: data.get('password'),
      });
      if (response.ok) {
        const text = await response.text();
        console.log('Raw response:', text);

        const rawResult = JSON.stringify(text)
        const result = JSON.parse(rawResult);

        console.log('login successful:', result);
        setMessage('login successful');
        setSeverity('success');

        const jsonParts = text.split('}')
                              .map(part => part.trim()) // 去除多余的空格
                              .filter(part => part) // 过滤掉空字符串
                              .map(part => part + '}'); // 再将分割的 "}" 加回去

        const tokenObject = JSON.parse(jsonParts[0]);
        const userInfo = JSON.parse(jsonParts[1]);
        const token = tokenObject.token;
        const userIdToSave = userInfo.id
        const userNameToSave = userInfo.username
        const userRoleToSave = userInfo.role
        console.log('Token:', token); // 这里会输出 token

        const expiresIn = 60 * 60 * 24 * 7; // Token expires in 7 days
        document.cookie = `token=${token}; max-age=${expiresIn}; path=/`;
        document.cookie = `savedUserId=${userIdToSave}; max-age=${expiresIn}; path=/`;
        document.cookie = `savedUserName=${userNameToSave}; max-age=${expiresIn}; path=/`;
        document.cookie = `savedUserRole=${userRoleToSave}; token=${token}; max-age=${expiresIn}; path=/`;

        window.location.href = '/dashboard';
      } else {
        console.log('in else');
        throw new Error('login failed');
        
        
      }
    } catch (error) {
      console.error('Error during login:', error);
      setMessage('login failed');
      setSeverity('error');
    }
    setOpen(true);
    
  };

  const handleClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
    if (reason === 'clickaway') {
      return;
    }

    setOpen(false);
  };

  return (
    <ThemeProvider theme={defaultTheme}>
      <Container component="main" maxWidth="xs">
        <CssBaseline />
        <Box
          sx={{
            marginTop: 8,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
          }}
        >
          <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }}>
            <LockOutlinedIcon />
          </Avatar>
          <Typography component="h1" variant="h5">
            Sign in
          </Typography>
          <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
            <TextField
              margin="normal"
              required
              fullWidth
              id="username"
              label="User Name"
              name="username"
              autoComplete="username"
              autoFocus
            />
            <TextField
              margin="normal"
              required
              fullWidth
              name="password"
              label="Password"
              type="password"
              id="password"
              autoComplete="current-password"
            />
            <FormControlLabel
              control={<Checkbox value="remember" color="primary" />}
              label="Remember me"
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
            >
              Sign In
            </Button>
            <Grid container>
              <Grid item xs>
                <Link href="#" variant="body2">
                  Forgot password?
                </Link>
              </Grid>
              <Grid item>
                <Link href="/sign-up" variant="body2">
                  {"Don't have an account? Sign Up"}
                </Link>
              </Grid>
            </Grid>
          </Box>
        </Box>
        <Copyright sx={{ mt: 8, mb: 4 }} />
      </Container>
    </ThemeProvider>
  );
}
