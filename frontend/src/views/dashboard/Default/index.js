import React from 'react';
import { useEffect, useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
// material-ui
import { Grid, Box, Button, TablePagination, Typography, Paper, Table, TableBody, TableContainer, TableHead, TableRow, TextField, MenuItem } from '@mui/material';
import { styled } from '@mui/material/styles';
import { Modal, ModalDialog } from '@mui/material';
import { useTranslation } from 'react-i18next';
// project imports
import { Dialog, DialogTitle, DialogContent, DialogActions, IconButton } from '@mui/material';
import axios from 'axios';
import SysUserInfoCard from './SysUserInfoCard';
import ConsumerInfoCard from './ConsumerInfoCard';
import ProducerInfoCard from './ProducerInfoCard';
import { DASHBOARD_URL } from '../../../UrlPath.js';
import { gridSpacing } from '../../../store/constant';
import Chart from "react-apexcharts";
import TableCell, { tableCellClasses } from '@mui/material/TableCell';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';

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

// ==============================|| DEFAULT DASHBOARD ||============================== //
// Modal content component
// Modal content component
const ModalContent = ({ consumer, onClose }) => (
  <div style={{ position: 'absolute', top: '50%', left: '50%', transform: 'translate(-50%, -50%)', width: 400, bgcolor: 'white', p: 4 }}>

    {/* Display service details inside the Modal */}
    <TableContainer component={Paper}>
      <Table size="small">
        <TableHead>
          <TableRow>
            <TableCell style={{ backgroundColor: '#FFD580', color: '#E35335', fontWeight: 'bold' }} colSpan={2}> {consumer.consumer_name} All Services Usage Details</TableCell>
          </TableRow>
          <TableRow>
            <TableCell style={{ backgroundColor: '#F0E68C', color: '#800000', fontWeight: 'bold' }}>Service</TableCell>
            <TableCell style={{ backgroundColor: '#F0E68C', color: '#800000', fontWeight: 'bold' }}>Total Usage  Count </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {consumer.services_data.map((service, index) => (
            <TableRow key={index}>
              <TableCell>{service.service_name}</TableCell>
              <TableCell>{service.service_count}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>

    {/* Manual close button */}
    <Button variant="contained" color="primary" onClick={onClose} style={{ position: 'absolute', top: '4px', right: '4px' }}>
      Close
    </Button>
  </div>
);

const Dashboard = () => {
  const [isLoading, setLoading] = useState(true);
  const [openPopup, setOpenPopup] = useState(false);
  const [popupMessage, setPopupMessage] = useState('');
  const [rows, setRows] = useState([]);
  const { t, i18n } = useTranslation();
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [page, setPage] = useState(0);
  const [expandedRows, setExpandedRows] = useState([]);
  const navigate = useNavigate();

  const today = dayjs();
  const [dateRange, setDateRange] = useState('today');
  const lastMonthStart = today.subtract(1, 'month').startOf('month');
  const lastMonthEnd = today.subtract(1, 'month').endOf('month');
  const [customDateRangeText, setCustomDateRangeText] = useState(
    `${lastMonthStart.format('MM/DD/YYYY')} - ${today.format('MM/DD/YYYY')}`
  );
  const [customDateRange, setCustomDateRange] = useState([
    lastMonthStart.format('MM/DD/YYYY'),
    today.format('MM/DD/YYYY'),
  ]);
  const [consumers, setConsumers] = useState([]);
  const [chartData, setChartData] = useState({  // Add this line
    series: [{
      data: [],
    }],
    options: {
      // ... (your other chart options)
    },
  });

  const dateRangePickerRef = useRef(null);

  useEffect(() => {
    setLoading(true);

    const fetchUserData = async () => {
      try {
        const response = await axios.get(DASHBOARD_URL);
        // console.log("Response from API:", response.data);

        if (response.data.SysUserPasswordSet === "false") {
          setPopupMessage("Please Reset Your System Password");
          setOpenPopup(true);
        }

        if (response.data && response.data.CustomerData) {
          setRows(response.data.CustomerData);
        } else {
          setRows([]);
        }

        const consumerData = response.data.ConsumerUsageData;
        // console.log("Consumer Data:", consumerData);
        setConsumers(consumerData);

        // Check if consumerData is null or undefined
        if (consumerData == null) {
          // Set default or empty data for the chart
          setChartData({
            series: [{
              data: [],
            }],
            options: {
              // ... (your other chart options)
            },
          });
        } else {

          // Set chartData only when data is loaded
          setChartData({
            series: [{
              data: consumerData.map((consumer) => getTotalCount(consumer.services_data)),
            }],
            options: {
              chart: {
                background: '#fff',
                height: 350,
                type: 'bar',
                events: {
                  click: (event, chartContext, config) => {
                    if (config.dataPointIndex !== undefined) {
                      const clickedConsumer = consumerData[config.dataPointIndex];
                      setSelectedConsumer(clickedConsumer);
                    }
                  },
                },
              },
              colors: consumerData.map((consumer) => consumer.consumer_color_code),
              plotOptions: {
                bar: {
                  columnWidth: '45%',
                  distributed: true,
                },
              },
              dataLabels: {
                enabled: true,
              },
              legend: {
                show: true,
              },
              xaxis: {
                categories: consumerData.map((consumer) => consumer.consumer_name),
                labels: {
                  style: {
                    colors: consumerData.map((consumer) => consumer.consumer_color_code),
                    fontSize: '12px',
                  },
                },
              },
              tooltip: {
                y: {
                  formatter: function (val, { seriesIndex, dataPointIndex, w }) {
                    const consumer = consumerData[dataPointIndex];
                    let tooltipText = `<div style="text-align: left; padding: 10px;">`;

                    // Add main label
                    tooltipText += `<span style="font-weight: bold; font-size: 12px;">${consumer.consumer_name} Service's Details</span>`;

                    // Add a list of service details
                    tooltipText += `<ul>`;
                    consumer.services_data.forEach((service) => {
                      tooltipText += `<li><b>${service.service_name}</b>: ${service.service_count}</li>`;
                    });
                    tooltipText += `</ul></div>`;

                    return tooltipText;
                  },
                },
              },

            },
          });


        }



        setLoading(false);
      } catch (error) {
        console.error('Error:', error);
        setLoading(false);
      }
    };

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
    // console.log('Date Range Picker is initialized');

    // Listen for date range changes
    $(dateRangePickerRef.current).on('apply.daterangepicker', (event, picker) => {
      const startDate = picker.startDate.format('MM/DD/YYYY');
      const endDate = picker.endDate.format('MM/DD/YYYY');

      // Log the selected date range
      // console.log(`Selected date range: ${startDate} - ${endDate}`);

      // Set the selected date range to the state
      setCustomDateRange([startDate, endDate]);

      // Update the input field with the selected date range
      setCustomDateRangeText(`${startDate} - ${endDate}`);
    });

    fetchUserData();
  }, []); // Add getTotalCount to the dependency array




  // const chartDataLoaded = consumers && consumers.length > 0;
  // Check if data has been loaded
  //console.log("chartDataLoaded",chartDataLoaded);




  const handleCloseDialog = () => {
    setOpenPopup(false);
    navigate('/ChangePassword'); // Redirect to /ChangePassword
  };

  const handleExpandClick = (rowId) => {
    if (expandedRows.includes(rowId)) {
      setExpandedRows(expandedRows.filter((id) => id !== rowId));
    } else {
      setExpandedRows([...expandedRows, rowId]);
    }
  };
  const isRowExpanded = (rowId) => expandedRows.includes(rowId);

  const StyledTableCell = styled(TableCell)(({ theme }) => ({
    [`&.${tableCellClasses.head}`]: {
      backgroundColor: theme.palette.common.black,
      color: theme.palette.common.white,
    },
    [`&.${tableCellClasses.body}`]: {
      fontSize: 18,
    },
  }));

  const StyledTableRow = styled(TableRow)(({ theme }) => ({
    '&:nth-of-type(odd)': {
      backgroundColor: theme.palette.action.hover,
    },
    // hide last border
    '&:last-child td, &:last-child th': {
      border: 0,
    },
  }));


  const getTotalCount = (servicesData) => {
    return servicesData.reduce((total, service) => total + parseInt(service.service_count, 10), 0);
  };
  // Add a state variable for the selected consumer
  const [selectedConsumer, setSelectedConsumer] = useState(null);

  //	------For daterange---------//




  const fetchDataByDateRange = async () => {
    try {
      const response = await axios.post(DASHBOARD_URL, {
        customStartDate: customDateRange[0],
        customEndDate: customDateRange[1],
      });
      const consumerData = response.data.ConsumerUsageDataPostMethod;
      //console.log("Consumer Data:", consumerData);
      setConsumers(consumerData);

      // Check if consumerData is null or undefined
      if (consumerData == null) {
        // Set default or empty data for the chart
        setChartData({
          series: [{
            data: [],
          }],
          options: {
            // ... (your other chart options)
          },
        });
      } else {

        // Set chartData only when data is loaded
        setChartData({
          series: [{
            data: consumerData.map((consumer) => getTotalCount(consumer.services_data)),
          }],
          options: {
            chart: {
              background: '#fff',
              height: 350,
              type: 'bar',
              events: {
                click: (event, chartContext, config) => {
                  if (config.dataPointIndex !== undefined) {
                    const clickedConsumer = consumerData[config.dataPointIndex];
                    setSelectedConsumer(clickedConsumer);
                  }
                },
              },
            },
            colors: consumerData.map((consumer) => consumer.consumer_color_code),
            plotOptions: {
              bar: {
                columnWidth: '45%',
                distributed: true,
              },
            },
            dataLabels: {
              enabled: true,
            },
            legend: {
              show: true,
            },
            xaxis: {
              categories: consumerData.map((consumer) => consumer.consumer_name),
              labels: {
                style: {
                  colors: consumerData.map((consumer) => consumer.consumer_color_code),
                  fontSize: '12px',
                },
              },
            },
            tooltip: {
              y: {
                formatter: function (val, { seriesIndex, dataPointIndex, w }) {
                  const consumer = consumerData[dataPointIndex];
                  let tooltipText = `<div style="text-align: left; padding: 10px;">`;

                  // Add main label
                  tooltipText += `<span style="font-weight: bold; font-size: 12px;">${consumer.consumer_name} Service's Details</span>`;

                  // Add a list of service details
                  tooltipText += `<ul>`;
                  consumer.services_data.forEach((service) => {
                    tooltipText += `<li><b>${service.service_name}</b>: ${service.service_count}</li>`;
                  });
                  tooltipText += `</ul></div>`;

                  return tooltipText;
                },
              },
            },

          },
        });

      }
      // Set chartData only when data is loaded



    } catch (error) {
      console.error('Error:', error);
      setLoading(false);
    }
  };





  return (
    <Grid container spacing={gridSpacing}>
      <Grid item xs={12}>
        <Grid container spacing={gridSpacing}>
          <Grid item lg={4} md={6} sm={6} xs={12}>
            <SysUserInfoCard isLoading={isLoading} />
          </Grid>
          <Grid item lg={4} md={6} sm={6} xs={12}>
            <ProducerInfoCard isLoading={isLoading} />
          </Grid>
          <Grid item lg={4} md={12} sm={12} xs={12}>
            <ConsumerInfoCard isLoading={isLoading} />
          </Grid>
        </Grid>
      </Grid>
      <Grid item lg={12} md={12} sm={12} xs={12}>
        <Grid container spacing={3} style={{ backgroundColor: 'white', borderRadius: '10px', margin: 'auto', width: '100%' }}>

          {/* Left side: h2 and h4 */}
          <Grid item lg={6} md={6} sm={6} xs={6}>
            <Grid container direction="column" spacing={1}>
              <Grid item>
                <Typography variant="h2">{t('ConsumerUsageStatistics')}</Typography>
              </Grid>
              <Grid item>
                <Typography variant="subtitle2" style={{ fontSize: '16px' }}>{t('ClickToGetConsumerUsageServiceDetails')}</Typography>
              </Grid>
            </Grid>
          </Grid>
          <Grid item lg={6} md={6} sm={6} xs={6} style={{ display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
            <div className="form-group" style={{ position: 'relative', marginRight: '20px' }}>
              <TextField
                label={t('customDateRange')}
                variant="outlined"
                margin="normal"
                value={customDateRangeText}
                InputProps={{
                  readOnly: true,
                  onClick: () => {
                    $(dateRangePickerRef.current).click();
                  },
                }}
              />
              {/* DateRangePicker will be initialized here */}
              <div ref={dateRangePickerRef} style={{ position: 'absolute', top: '100%', left: 0 }} />
            </div>
            <Button variant="contained" color="primary" style={{ marginRight: '20px' }} onClick={fetchDataByDateRange}>
              {t('searchButton')}
            </Button>
          </Grid>
          <Grid item lg={12} md={12} sm={12} xs={12}>
            {/* Render the Chart only if data is loaded */}
            {/* {chartDataLoaded && ( */}
            <Grid item lg={12} md={12} sm={12} xs={12}>
              {chartData.series[0].data && chartData.series[0].data.length > 0 ? (
                <Chart
                  options={chartData.options}
                  series={chartData.series}
                  type="bar"
                  height={350}
                />
              ) : (
                <div style={{ textAlign: 'center', marginTop: '20px' }}>
                  <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                    {sessionStorage.getItem('language') === 'english'
                      ? 'No Usage Data Available'
                      : 'Aucune donnée d\'utilisation disponible'}
                  </Typography><br />
                </div>
              )}
            </Grid>

            {/* )} */}
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12} style={{ backgroundColor: 'white', borderRadius: '10px', margin: 'auto', width: '100%', marginTop: '20px', marginLeft: '25px' }}>
        <Grid container spacing={gridSpacing}>
          <Grid item lg={12} md={12} sm={12} xs={12}>
            <>
              <Grid container direction="column" spacing={1}>
                <Grid item>
                  <Typography variant="h2">{t('ProducerToConsumerMappedServicesDetails')}</Typography>
                </Grid>
              </Grid>
              <br />
              <TableContainer component={Paper} style={{ marginLeft: '-10px' }}>
                <Table id="tableId" sx={{ minWidth: 700 }} aria-label="customized table">
                  <TableHead>
                    <TableRow>
                      <StyledTableCell >{t('producer')}</StyledTableCell>
                      <StyledTableCell >{t('consumer')}</StyledTableCell>
                      <StyledTableCell >{t('subscribedServices')}</StyledTableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {(rowsPerPage > 0
                      ? rows.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                      : rows
                    ).map((row) => (
                      <React.Fragment key={row.Id}>
                        <StyledTableRow>
                          <StyledTableCell>{row.ProducerName}</StyledTableCell>
                          <StyledTableCell>{row.ConsumerName}</StyledTableCell>
                          <StyledTableCell>
                            <IconButton
                              aria-label="expand row"
                              size="small"
                              onClick={() => handleExpandClick(row.Id)}
                            >
                              {isRowExpanded(row.Id) ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
                            </IconButton>
                          </StyledTableCell>
                        </StyledTableRow>
                        {isRowExpanded(row.Id) && (
                          <TableRow >
                            <TableCell colSpan={3}>
                              <Table size="small">
                                <TableHead>
                                  <TableRow>
                                    <TableCell sx={{ fontWeight: 'bold', background: '#814141', color: 'white', }}>Service Name</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold', background: '#814141', color: 'white', }}>Service URL</TableCell>
                                  </TableRow>
                                </TableHead>
                                <TableBody>
                                  {JSON.parse(row.ProducerServices).subsrcibed_services.map((service, index) => (
                                    <TableRow key={index}>
                                      <TableCell>{service.service_name}</TableCell>
                                      <TableCell>{row.ConsumerDomainAddress}{service.service_url}</TableCell>
                                    </TableRow>
                                  ))}
                                </TableBody>
                              </Table>
                            </TableCell>
                          </TableRow>
                        )}
                      </React.Fragment>
                    ))}
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
            </>

          </Grid>


          {/* Conditionally render the Modal */}
          {selectedConsumer && (
            <Modal open={Boolean(selectedConsumer)} onClose={() => setSelectedConsumer(null)}>
              <ModalContent consumer={selectedConsumer} onClose={() => setSelectedConsumer(null)} />
            </Modal>
          )}
        </Grid>
      </Grid>
      {/* Dialog */}
      <Dialog open={openPopup} onClose={handleCloseDialog}>
        <DialogTitle>SuperESB Admin Alert</DialogTitle>
        <DialogContent>
          {popupMessage}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDialog} color="primary">
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </Grid>
  );
};

export default Dashboard;