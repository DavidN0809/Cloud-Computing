import React, { useEffect, useState } from 'react';
import { List, ListItem, ListItemText, Paper, Typography, Button } from '@mui/material';

const ListAllBillings = ({ setBillingId }) => {
  const [billings, setBillings] = useState([]);

  useEffect(() => {
    const fetchBillings = async () => {
      const response = await fetch('/billings/list');
      const data = await response.json();
      if (response.ok) {
        setBillings(data);
      } else {
        alert('Failed to fetch billings');
      }
    };

    fetchBillings();
  }, []);

  const handleDelete = async (billingId) => {
    if (window.confirm("Are you sure you want to delete this billing record?")) {
      const response = await fetch(`/billings/remove/${billingId}`, {
        method: 'DELETE'
      });
      if (response.ok) {
        setBillings(billings.filter(billing => billing.ID !== billingId));
        alert("Billing record deleted successfully.");
      } else {
        alert("Failed to delete billing record.");
      }
    }
  };

  return (
    <Paper sx={{ padding: 2, marginTop: 2 }}>
      <Typography variant="h6">All Billings</Typography>
      <List>
        {billings.map((billing) => (
          <ListItem key={billing.ID} secondaryAction={
            <>
              <Button onClick={() => setBillingId(billing.ID)} sx={{ marginRight: 2 }}>Edit</Button>
              <Button onClick={() => handleDelete(billing.ID)} color="error">Delete</Button>
            </>
          }>
            <ListItemText
              primary={`Billing ID: ${billing.ID}`}
              secondary={`Hours: ${billing.hours}, Amount: $${billing.amount}, Task ID: ${billing.taskID}, User ID: ${billing.userID}`}
            />
          </ListItem>
        ))}
      </List>
    </Paper>
  );
};

export default ListAllBillings;

