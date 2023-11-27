import React, { useState, useEffect, useRef } from 'react';
import { styled } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';
import { TablePagination } from '@mui/material';
import { useTranslation } from 'react-i18next';
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Typography } from '@mui/material';
import { Card, Grid, Box, Button, TextField, FormControl, FormHelperText, InputLabel, Select, MenuItem, IconButton } from '@mui/material';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import InsertDriveFileIcon from '@mui/icons-material/InsertDriveFile';
import PictureAsPdfIcon from '@mui/icons-material/PictureAsPdf';




// Import the daterangepicker CSS
import dayjs from 'dayjs';
import 'dayjs/locale/en';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-daterangepicker/daterangepicker.css';
import moment from 'moment';
import $ from 'jquery';
import 'popper.js';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import 'bootstrap-daterangepicker/daterangepicker.js';
import CopyToClipboard from 'react-copy-to-clipboard';
import * as XLSX from 'xlsx';
import jsPDF from 'jspdf';
import 'jspdf-autotable';

//for making HTTP requests
import axios from 'axios';

// project imports
import MainCard from '../../../ui-component/cards/MainCard';
import SecondaryAction from '../../../ui-component/cards/CardSecondaryAction';
import { CREATE_PRODUCERS_URL, SEARCH_PRODUCERS_URL, ESBLOGS_REPORT_URL } from '../../../UrlPath.js';

// assets
import LinkIcon from '@mui/icons-material/Link';
import '../../../assets/scss/style.css'

// styles
const IFrameWrapper = styled('iframe')(({ theme }) => ({
  height: 'calc(100vh - 210px)',
  border: '1px solid',
  borderColor: theme.palette.primary.light
}));

// =============================|| TABLER ICONS ||============================= //

const ConsumerReports = () => {
  const navigate = useNavigate();
  const [errors, setErrors] = useState({});
  const [rows, setRows] = useState([]);
  const { t, i18n } = useTranslation();
  const [customDateRangeText, setCustomDateRangeText] = useState('');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);

// Data to copy or export to xlsx or pdf in table
    const [copied, setCopied] = useState(false);
    const [tableData, setTableData] = useState([]);
    
   // Add the exportToExcel function
	const exportToExcel = () => {
  // Customize the data before exporting
  const customizedTableData = tableData.map((row) => ({
	'Id':row.Id,
    'Full Name': row.FullName,
    'Email': row.Email,
    'Status': row.Status,
    'Timestamp':row.Timestamp,
    'Mobile Number':row.Mobile,
    'Address':row.Address,
    
  }));

  const ws = XLSX.utils.json_to_sheet(customizedTableData);
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, 'Sheet1');
  XLSX.writeFile(wb, 'ConsumerReportData.xlsx');
};

// Add the exportToPdf function
const exportToPDF = () => {
  const doc = new jsPDF();
  doc.text('Table Data', 10, 10);
  
  // Check the table ID used in your HTML, it should match your table's ID
  doc.autoTable({ startY: 20, html: '#tableId' });

  doc.save('ConsumerReportData.pdf');
};




  //	------For daterange---------//

  const dateRangePickerRef = useRef(null);

  useEffect(() => {
    $(document).ready(function () {
      // Log a message to check if this block is reached
      console.log('Document is ready');

      // Customize the Date Range Picker options
      const dateRangePickerOptions = {
    opens: 'left',
    autoUpdateInput: false,
    locale: {
      format: 'MM/DD/YYYY',
      applyLabel: t('Apply'),
      cancelLabel: t('Cancel'),
      customRangeLabel: t('Custom Range'),
    },
    ranges: {
      [t('Today')]: [moment(), moment()],
      [t('Yesterday')]: [moment().subtract(1, 'days'), moment().subtract(1, 'days')],
      [t('Last 7 Days')]: [moment().subtract(6, 'days'), moment()],
      [t('Last 30 Days')]: [moment().subtract(29, 'days'), moment()],
      [t('This Month')]: [moment().startOf('month'), moment().endOf('month')],
      [t('Last Month')]: [
        moment().subtract(1, 'month').startOf('month'),
        moment().subtract(1, 'month').endOf('month'),
      ],
        },
      };

      // Initialize the date range picker with options
      $(dateRangePickerRef.current).daterangepicker(dateRangePickerOptions);

      // Log a message to check if the Date Range Picker is initialized
      console.log('Date Range Picker is initialized');

      // Listen for date range changes
      $(dateRangePickerRef.current).on('apply.daterangepicker', (event, picker) => {
        const startDate = picker.startDate.format('MM/DD/YYYY');
        const endDate = picker.endDate.format('MM/DD/YYYY');

        // Log the selected date range
        console.log(`Selected date range: ${startDate} - ${endDate}`);

        // Set the selected date range to the state
        setCustomDateRange([startDate, endDate]);

        // Update the input field with the selected date range
        setCustomDateRangeText(`${startDate} - ${endDate}`);
      });
    });
  }, []);

  const today = dayjs();
  const lastMonthStart = today.subtract(1, 'month').startOf('month');
  const lastMonthEnd = today.subtract(1, 'month').endOf('month');

  const [customDateRange, setCustomDateRange] = useState([lastMonthStart, today]);

  //  --------For status dropdown------

  const hardcodedStatuses = [
    { id: 'ACTIVE', name: 'Active' },
    { id: 'INACTIVE', name: 'Inactive' },
  ];

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(SUBSCRIBER_REPORT_URL);
         setTableData(response.data.esb_logs_data); 
        setRows(response.data.esb_logs_data);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, []);


  const [formData, setFormData] = useState({
    roleName: '',
    input_status: '',
  });


  const handleFormSubmit = async () => {
    try {
      const response = await axios.post(SEARCH_SUBSCRIBER_URL, {
        ...formData,
        customStartDate: customDateRange[0],
        customEndDate: customDateRange[1],
      });

      if (response.data && response.data.CustomerDataPostMethod) {
        if (response.data.CustomerDataPostMethod.length === 0) {
          setRows([]); // Set rows to an empty array to trigger "No data available" message
        } else {
          setRows(response.data.CustomerDataPostMethod); // Set rows with the response data
        }
      } else {
        setRows([]); // Set rows to an empty array if the data is undefined or missing
      }
    } catch (error) {
      // Handle network or other errors
      console.error('Error searching users:', error);
      setRows([]); // Set rows to an empty array on error to show "No data available"
    }
  };

  const handleChange = (name, value) => {
    setFormData((prevFormData) => ({
      ...prevFormData,
      [name]: value,
    }));
  };

 


  return (
    <>
      <MainCard
        title={<span style={{ color: 'black ', fontWeight: 'bold', fontSize: '24px' }}>{t('ConsumerReport')}</span>}
      //  style={{ backgroundColor: '#4682B4' }}
        secondary={<SecondaryAction icon={<LinkIcon fontSize="small" />} link="https://tablericons.com/" />}
      >
        <Card sx={{ overflow: 'hidden' }} >
          <Grid container spacing={2} style={{ marginLeft: '10px' }}>
            <Grid item xs={3}>
              <TextField
                label={t('consumer')}
                variant="outlined"
                fullWidth
                margin="normal"
                value={formData.fullName}
                onChange={(e) => setFormData({ ...formData, fullName: e.target.value })}
              />
            </Grid>
            <Grid item xs={3}>
              <TextField
                label={t('RequestID')}
                variant="outlined"
                fullWidth
                margin="normal"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              />
            </Grid>
            <Grid item xs={3}>
              <TextField
                label={t('consumerCode')}
                variant="outlined"
                fullWidth
                margin="normal"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              />
            </Grid>
            {/* <Grid item xs={3}>
              <FormControl fullWidth variant="outlined" margin="normal">
                <InputLabel id="status-label" className={errors.language ? 'input-label-error' : ''}>
                  Status
                </InputLabel>
                <Select
                  labelId="status-label"
                  id="status"
                  value={formData.input_status}
                  onChange={(e) => handleChange('input_status', e.target.value)}
                  label="Status"
                  error={!!errors.status}
                >
                  <MenuItem value="">Select Status</MenuItem>
                  {hardcodedStatuses.map((status) => (
                    <MenuItem key={status.id} value={status.id}>
                      {status.name}
                    </MenuItem>
                  ))}
                </Select>
                {errors.status && <FormHelperText error>{errors.status}</FormHelperText>}
              </FormControl>
            </Grid> */}
            <Grid item xs={3}>
              <div className="form-group" style={{ position: 'relative' }}>
                <TextField
                  label={t('customDateRange')}
                  variant="outlined"
                  margin="normal"
                  value={customDateRangeText}
                  InputProps={{
                    readOnly: true,
                    onClick: () => {
                      // Log a message to check if the click event is triggered
                      console.log('Text field clicked');

                      // Show the date range picker when the text field is clicked
                      $(dateRangePickerRef.current).click();
                    },
                  }}
                />
                {/* DateRangePicker will be initialized here */}
                <div ref={dateRangePickerRef} style={{ position: 'absolute', top: '100%', left: 0 }} />
              </div>
            </Grid>
          </Grid>
          <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', p: 2 }}>
            <Button variant="contained" color="secondary" onClick={handleFormSubmit}>
              {t('searchButton')}
            </Button>
            <Box sx={{ mx: 2 }} />
          
          </Box>
        </Card>
      </MainCard>
      <br />
      <div>
       {/* copy button */}
       <Button variant="contained" startIcon={<ContentCopyIcon />} className="exportbutton" onClick={() => setCopied(false)}>
              <CopyToClipboard text={JSON.stringify(tableData)} onCopy={() => setCopied(true)}>
                <span>{t('Copy')}</span>
              </CopyToClipboard>
            </Button>
            {copied ? <span style={{ color: 'green' }}>Copied!</span> : null}
            {/* copy button */}

            {/* excel button */}
            <Button variant="contained" startIcon={<InsertDriveFileIcon />} className="exportbutton" onClick={exportToExcel}>
              {t('Excel')}
            </Button>

            {/* pdf button */}
            <Button variant="contained" startIcon={<PictureAsPdfIcon />} className="exportbutton" onClick={exportToPDF}>
              {t('Pdf')}
            </Button>
    </div>
      <MainCard
        title=""
        style={{ backgroundColor: '#CCCCFF' }}
        secondary={<SecondaryAction icon={<LinkIcon fontSize="small" />} link="https://tablericons.com/" />}
      >
        <Card sx={{ overflow: 'hidden' }}>{/* Your content here */}</Card>
        
        <TableContainer component={Paper}>
        
          <Table id="tableId">
            <TableHead style={{ backgroundColor: '#1434A4' }}>
              <TableRow>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('idLabel')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('timestampLabel')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('RequestID')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('consumer')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('consumerCode')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('InRequest')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('OutResponse')}</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
            <TableRow>
                <TableCell>19d9b4ff-4deb-4dc3-816e-35f69eb148bd</TableCell>
                <TableCell>2023-09-13 08:30:12</TableCell>
                <TableCell>1944737</TableCell>
                <TableCell>ECNPS-Portal</TableCell>
                <TableCell>dfw8943847jhduy4328</TableCell>
                <TableCell>"producer":"ecnps"</TableCell>
                <TableCell>"details":"pgs1 api is called - e-cnps details sent"</TableCell>
              </TableRow>
              <TableRow>
                <TableCell>19d9b4ff-4deb-4dc3-816e-35f69eb148bd</TableCell>
                <TableCell>2023-09-13 08:30:12</TableCell>
                <TableCell>1944737</TableCell>
                <TableCell>ECNPS-Portal</TableCell>
                <TableCell>dfw8943847jhduy4328</TableCell>
                <TableCell>"producer":"ecnps"</TableCell>
                <TableCell>"details":"pgs1 api is called - e-cnps details sent"</TableCell>
              </TableRow>
              <TableRow>
                <TableCell>19d9b4ff-4deb-4dc3-816e-35f69eb148bd</TableCell>
                <TableCell>2023-09-13 08:30:12</TableCell>
                <TableCell>1944737</TableCell>
                <TableCell>ECNPS-Portal</TableCell>
                <TableCell>dfw8943847jhduy4328</TableCell>
                <TableCell>"producer":"ecnps"</TableCell>
                <TableCell>"details":"pgs1 api is called - e-cnps details sent"</TableCell>
              </TableRow>
              <TableRow>
                <TableCell>19d9b4ff-4deb-4dc3-816e-35f69eb148bd</TableCell>
                <TableCell>2023-09-13 08:30:12</TableCell>
                <TableCell>1944737</TableCell>
                <TableCell>ECNPS-Portal</TableCell>
                <TableCell>dfw8943847jhduy4328</TableCell>
                <TableCell>"producer":"ecnps"</TableCell>
                <TableCell>"details":"pgs1 api is called - e-cnps details sent"</TableCell>
              </TableRow>
              <TableRow>
                <TableCell>19d9b4ff-4deb-4dc3-816e-35f69eb148bd</TableCell>
                <TableCell>2023-09-13 08:30:12</TableCell>
                <TableCell>1944737</TableCell>
                <TableCell>ECNPS-Portal</TableCell>
                <TableCell>dfw8943847jhduy4328</TableCell>
                <TableCell>"producer":"ecnps"</TableCell>
                <TableCell>"details":"pgs1 api is called - e-cnps details sent"</TableCell>
              </TableRow>
              <TableRow>
                <TableCell>19d9b4ff-4deb-4dc3-816e-35f69eb148bd</TableCell>
                <TableCell>2023-09-13 08:30:12</TableCell>
                <TableCell>1944737</TableCell>
                <TableCell>ECNPS-Portal</TableCell>
                <TableCell>dfw8943847jhduy4328</TableCell>
                <TableCell>"producer":"ecnps"</TableCell>
                <TableCell>"details":"pgs1 api is called - e-cnps details sent"</TableCell>
              </TableRow>
              <TableRow>
                <TableCell>19d9b4ff-4deb-4dc3-816e-35f69eb148bd</TableCell>
                <TableCell>2023-09-13 08:30:12</TableCell>
                <TableCell>1944737</TableCell>
                <TableCell>ECNPS-Portal</TableCell>
                <TableCell>dfw8943847jhduy4328</TableCell>
                <TableCell>"producer":"ecnps"</TableCell>
                <TableCell>"details":"pgs1 api is called - e-cnps details sent"</TableCell>
              </TableRow>
              <TableRow>
                <TableCell>19d9b4ff-4deb-4dc3-816e-35f69eb148bd</TableCell>
                <TableCell>2023-09-13 08:30:12</TableCell>
                <TableCell>1944737</TableCell>
                <TableCell>ECNPS-Portal</TableCell>
                <TableCell>dfw8943847jhduy4328</TableCell>
                <TableCell>"producer":"ecnps"</TableCell>
                <TableCell>"details":"pgs1 api is called - e-cnps details sent"</TableCell>
              </TableRow>
              <TableRow>
                <TableCell>19d9b4ff-4deb-4dc3-816e-35f69eb148bd</TableCell>
                <TableCell>2023-09-13 08:30:12</TableCell>
                <TableCell>1944737</TableCell>
                <TableCell>ECNPS-Portal</TableCell>
                <TableCell>dfw8943847jhduy4328</TableCell>
                <TableCell>"producer":"ecnps"</TableCell>
                <TableCell>"details":"pgs1 api is called - e-cnps details sent"</TableCell>
              </TableRow>
              <TableRow>
                <TableCell>19d9b4ff-4deb-4dc3-816e-35f69eb148bd</TableCell>
                <TableCell>2023-09-13 08:30:12</TableCell>
                <TableCell>1944737</TableCell>
                <TableCell>ECNPS-Portal</TableCell>
                <TableCell>dfw8943847jhduy4328</TableCell>
                <TableCell>"producer":"ecnps"</TableCell>
                <TableCell>"details":"pgs1 api is called - e-cnps details sent"</TableCell>
              </TableRow>
              <TableRow>
                <TableCell>19d9b4ff-4deb-4dc3-816e-35f69eb148bd</TableCell>
                <TableCell>2023-09-13 08:30:12</TableCell>
                <TableCell>1944737</TableCell>
                <TableCell>ECNPS-Portal</TableCell>
                <TableCell>dfw8943847jhduy4328</TableCell>
                <TableCell>"producer":"ecnps"</TableCell>
                <TableCell>"details":"pgs1 api is called - e-cnps details sent"</TableCell>
              </TableRow>
            </TableBody>

          </Table>

          <TablePagination
            rowsPerPageOptions={[5, 10, 25, 50]}
            component="div"
            count={rows.length}
            rowsPerPage={rowsPerPage}
            page={page}
            onPageChange={(event, newPage) => setPage(newPage)}
            onRowsPerPageChange={(event) => {
              setRowsPerPage(parseInt(event.target.value, 10));
              setPage(0);
            }}
            labelRowsPerPage={t('rowsPerPageLabel')}
            labelDisplayedRows={({ from, to, count }) => `${from}-${to} of ${count}`}
            // Custom style for the pagination controls
            className="pagination-container"
          />
        </TableContainer>
        
      </MainCard>
    </>
  );
};

export default ConsumerReports;
