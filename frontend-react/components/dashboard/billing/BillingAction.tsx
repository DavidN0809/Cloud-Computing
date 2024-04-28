import React, { useState, useEffect } from 'react';
import { Button, TextField, Box, Typography } from '@mui/material';

const BillingAction = ({ billingId, setBillingId }) => {
  const [billingDetails, setBillingDetails] = useState({
    hours: '',
    amount: '',
    taskID: '',
    userID: ''
  });

  useEffect(() => {
    if (billingId) {
      fetch(`/billings/get/${billingId}`)
        .then(res => res.json())
        .then(data => {
          setBillingDetails(data);
        });
    }
  }, [billingId]);

  const handleChange = (event) => {
    setBillingDetails({
      ...billingDetails,
      [event.target.name]: event.target.value
    });
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    const method = billingId ? 'PUT' : 'POST';
    const url = billingId ? `/billings/update/${billingId}` : '/billings/create';

    const response = await fetch(url, {
      method: method,
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(billingDetails)
    });
    const data = await response.json();
    if (response.ok) {
      alert('Billing operation successful');
      setBillingDetails({
        hours: '',
        amount: '',
        taskID: '',
        userID: ''
      }); // Reset form
      setBillingId(null); // Reset billingId
    } else {
      alert('Failed to perform billing operation: ' + data.message);
    }
  };

  return (
    <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
      <Typography variant="h6">{billingId ? "Edit Billing" : "Create Billing"}</Typography>
      <TextField
        margin="normal"
        required
        fullWidth
        name="hours"
        label="Hours"
        type="number"
        id="hours"
        value={billingDetails.hours}
        onChange={handleChange}
      />
      <TextField
        margin="normal"
        required
        fullWidth
        name="amount"
        label="Amount"
        type="number"
        id="amount"
        value={billingDetails.amount}
        onChange={handleChange}
      />
      <TextField
        margin="normal"
        required
        fullWidth
        name="taskID"
        label="Task ID"
        id="taskID"
        value={billingDetails.taskID}
        onChange={handleChange}
      />
      <TextField
        margin="normal"
        required
        fullWidth
        name="userID"
        label="User ID"
        id="userID"
        value={billingDetails.userID}
        onChange={handleChange}
      />
      <Button
        type="submit"
        fullWidth
        variant="contained"
        sx={{ mt: 3, mb: 2 }}
      >
        {billingId ? "Update" : "Create"}
      </Button>
    </Box>
  );
};

export default BillingAction;

