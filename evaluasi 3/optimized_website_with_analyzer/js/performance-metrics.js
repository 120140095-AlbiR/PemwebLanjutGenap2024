/**
 * Performance Measurement System
 * 
 * This script tracks 5 key performance metrics:
 * 1. Page Load Time
 * 2. First Contentful Paint (FCP)
 * 3. Largest Contentful Paint (LCP)
 * 4. Cumulative Layout Shift (CLS)
 * 5. First Input Delay (FID)
 */

class PerformanceMetrics {
    constructor() {
        // Initialize metrics storage
        this.metrics = {
            pageLoadTime: 0,
            fcp: 0,
            lcp: 0,
            cls: 0,
            fid: 0,
        };

        this.clsValue = 0;
        this.clsEntries = [];

        // Dashboard element
        this.dashboardElement = null;

        // Initialize measurement
        this.initMeasurements();
    }

    /**
     * Initialize all performance measurements
     */
    initMeasurements() {
        // Wait for page to load to ensure we can measure everything
        window.addEventListener('load', () => {
            // Measure page load time
            this.measurePageLoadTime();

            // Measure FCP
            this.measureFCP();

            // Measure LCP
            this.measureLCP();

            // Measure CLS
            this.measureCLS();

            // Measure FID
            this.measureFID();

            // Create dashboard after a short delay to ensure metrics are collected
            setTimeout(() => {
                this.createDashboard();
            }, 3000);

            // Send analytics data
            setTimeout(() => {
                this.sendAnalytics();
            }, 5000);
        });
    }

    /**
     * Measure total page load time
     */
    measurePageLoadTime() {
        if (window.performance && window.performance.timing) {
            const timing = window.performance.timing;
            const loadTime = timing.loadEventEnd - timing.navigationStart;
            this.metrics.pageLoadTime = loadTime;
            console.log(`Page Load Time: ${loadTime}ms`);
        } else {
            // For newer browsers using Navigation Timing API Level 2
            const pageNav = performance.getEntriesByType('navigation')[0];
            if (pageNav) {
                this.metrics.pageLoadTime = pageNav.loadEventEnd;
                console.log(`Page Load Time: ${pageNav.loadEventEnd}ms`);
            }
        }
    }

    /**
     * Measure First Contentful Paint
     */
    measureFCP() {
        const fcpObserver = new PerformanceObserver((entryList) => {
            const entries = entryList.getEntries();
            const fcpEntry = entries[entries.length - 1];
            this.metrics.fcp = fcpEntry.startTime;
            console.log(`FCP: ${fcpEntry.startTime}ms`);
            fcpObserver.disconnect();
        });

        // Register observer
        fcpObserver.observe({ type: 'paint', buffered: true });
    }

    /**
     * Measure Largest Contentful Paint
     */
    measureLCP() {
        const lcpObserver = new PerformanceObserver((entryList) => {
            const entries = entryList.getEntries();
            const lcpEntry = entries[entries.length - 1];
            this.metrics.lcp = lcpEntry.startTime;
            console.log(`LCP: ${lcpEntry.startTime}ms`);
        });

        // Register observer
        lcpObserver.observe({ type: 'largest-contentful-paint', buffered: true });

        // Disconnect on visibility change
        document.addEventListener('visibilitychange', () => {
            if (document.visibilityState === 'hidden') {
                lcpObserver.disconnect();
            }
        });
    }

    /**
     * Measure Cumulative Layout Shift
     */
    measureCLS() {
        let clsValue = 0;
        let clsEntries = [];

        const clsObserver = new PerformanceObserver((entryList) => {
            for (const entry of entryList.getEntries()) {
                // Only count layout shifts without recent user input
                if (!entry.hadRecentInput) {
                    clsValue += entry.value;
                    clsEntries.push(entry);
                }
            }
            
            this.clsValue = clsValue;
            this.clsEntries = clsEntries;
            this.metrics.cls = clsValue;
            console.log(`Current CLS: ${clsValue}`);
        });

        // Register observer
        clsObserver.observe({ type: 'layout-shift', buffered: true });

        // Report final CLS when page is hidden
        document.addEventListener('visibilitychange', () => {
            if (document.visibilityState === 'hidden') {
                clsObserver.disconnect();
                console.log(`Final CLS: ${clsValue}`);
            }
        });
    }

    /**
     * Measure First Input Delay
     */
    measureFID() {
        const fidObserver = new PerformanceObserver((entryList) => {
            for (const entry of entryList.getEntries()) {
                this.metrics.fid = entry.processingStart - entry.startTime;
                console.log(`FID: ${this.metrics.fid}ms`);
                fidObserver.disconnect();
            }
        });

        // Register observer
        fidObserver.observe({ type: 'first-input', buffered: true });
    }

    /**
     * Create a visual dashboard to display performance metrics
     */
    createDashboard() {
        // Create dashboard container
        const dashboard = document.createElement('div');
        dashboard.id = 'performance-dashboard';
        dashboard.style.position = 'fixed';
        dashboard.style.bottom = '20px';
        dashboard.style.right = '20px';
        dashboard.style.backgroundColor = 'rgba(0, 0, 0, 0.8)';
        dashboard.style.color = 'white';
        dashboard.style.padding = '10px';
        dashboard.style.borderRadius = '5px';
        dashboard.style.zIndex = '9999';
        dashboard.style.fontSize = '12px';
        dashboard.style.boxShadow = '0 0 10px rgba(0, 0, 0, 0.5)';
        dashboard.style.maxWidth = '300px';
        dashboard.style.fontFamily = 'Arial, sans-serif';

        // Add a title
        const title = document.createElement('h3');
        title.textContent = 'Performance Metrics';
        title.style.margin = '0 0 10px 0';
        title.style.padding = '0 0 5px 0';
        title.style.borderBottom = '1px solid #666';
        title.style.fontSize = '14px';
        dashboard.appendChild(title);

        // Create a table for metrics
        const table = document.createElement('table');
        table.style.width = '100%';
        table.style.borderCollapse = 'collapse';
        
        // Add metrics to table
        this.addMetricRow(table, 'Page Load', `${this.metrics.pageLoadTime}ms`, this.getRatingColor(this.metrics.pageLoadTime, 3000));
        this.addMetricRow(table, 'FCP', `${this.metrics.fcp}ms`, this.getRatingColor(this.metrics.fcp, 1800));
        this.addMetricRow(table, 'LCP', `${this.metrics.lcp}ms`, this.getRatingColor(this.metrics.lcp, 2500));
        this.addMetricRow(table, 'CLS', this.metrics.cls.toFixed(3), this.getRatingColor(this.metrics.cls * 1000, 100, true));
        this.addMetricRow(table, 'FID', `${this.metrics.fid}ms`, this.getRatingColor(this.metrics.fid, 100));

        dashboard.appendChild(table);

        // Add close button
        const closeButton = document.createElement('button');
        closeButton.textContent = 'Close';
        closeButton.style.marginTop = '10px';
        closeButton.style.padding = '3px 8px';
        closeButton.style.backgroundColor = '#444';
        closeButton.style.border = 'none';
        closeButton.style.color = 'white';
        closeButton.style.borderRadius = '3px';
        closeButton.style.cursor = 'pointer';
        closeButton.onclick = () => {
            dashboard.style.display = 'none';
        };
        dashboard.appendChild(closeButton);

        // Append dashboard to body
        document.body.appendChild(dashboard);
        this.dashboardElement = dashboard;
    }

    /**
     * Add a row to the metrics table
     */
    addMetricRow(table, name, value, color) {
        const row = document.createElement('tr');
        
        const nameCell = document.createElement('td');
        nameCell.textContent = name;
        nameCell.style.padding = '3px 0';
        row.appendChild(nameCell);
        
        const valueCell = document.createElement('td');
        valueCell.textContent = value;
        valueCell.style.textAlign = 'right';
        valueCell.style.padding = '3px 0';
        valueCell.style.color = color;
        row.appendChild(valueCell);
        
        table.appendChild(row);
    }

    /**
     * Get a color rating based on metric value
     */
    getRatingColor(value, threshold, reverse = false) {
        if (reverse) {
            if (value < threshold * 0.33) return '#4caf50'; // Good
            if (value < threshold * 0.66) return '#ff9800'; // Needs Improvement
            return '#f44336'; // Poor
        } else {
            if (value < threshold * 0.33) return '#4caf50'; // Good
            if (value < threshold) return '#ff9800'; // Needs Improvement
            return '#f44336'; // Poor
        }
    }

    /**
     * Send analytics data to server (mock implementation)
     */
    sendAnalytics() {
        console.log('Sending performance analytics:', this.metrics);
        // In a real implementation, you would send this data to your analytics server
        // Example:
        // fetch('/api/analytics/performance', {
        //   method: 'POST',
        //   body: JSON.stringify(this.metrics),
        //   headers: { 'Content-Type': 'application/json' }
        // });
    }
}

// Initialize performance tracking when the script loads
const performanceTracker = new PerformanceMetrics();

/**
 * Improvement recommendations based on common issues
 * 
 * 1. Page Load Time:
 *    - Minimize and combine CSS/JS files
 *    - Enable browser caching
 *    - Optimize images
 *    - Use a CDN
 *
 * 2. First Contentful Paint (FCP):
 *    - Eliminate render-blocking resources
 *    - Minimize CSS
 *    - Remove unused CSS
 *
 * 3. Largest Contentful Paint (LCP):
 *    - Optimize images
 *    - Preload important resources
 *    - Implement server-side rendering
 *    - Use a CDN
 *
 * 4. Cumulative Layout Shift (CLS):
 *    - Set explicit width and height for images and videos
 *    - Avoid inserting content above existing content
 *    - Use transform animations instead of animations that trigger layout changes
 *
 * 5. First Input Delay (FID):
 *    - Break up Long Tasks
 *    - Optimize JavaScript execution
 *    - Minimize main thread work
 *    - Keep request counts low and transfer sizes small
 */
