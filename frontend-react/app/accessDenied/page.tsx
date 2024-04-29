import React from 'react';
import { Container, Typography, Button, Box } from '@mui/material';
import Head from 'next/head';
import Link from 'next/link';

const AccessDeniedPage = () => {
  return (
    <>
      <Head>
        <title>Access Denied</title>
      </Head>
      <Container component="main" maxWidth="sm" sx={{ pt: 8, pb: 6 }}>
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            minHeight: '80vh', // take up at least 80% of the viewport height
          }}
        >
          <Typography component="h1" variant="h2" align="center" gutterBottom>
            Access Denied
          </Typography>
          <Typography variant="h5" align="center">
            Sorry, you don't have access to this page.
          </Typography>
          <Link href="/dashboard" passHref>
            <Button variant="outlined" sx={{ mt: 3 }}>
              Go Back Home
            </Button>
          </Link>
        </Box>
      </Container>
    </>
  );
};

export default AccessDeniedPage;
