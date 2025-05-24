/**
 * SEO Analysis Tool
 * 
 * This script analyzes the current page for SEO best practices and provides recommendations.
 * It checks:
 * 1. Title tag optimization
 * 2. Meta description 
 * 3. Heading structure
 * 4. Image alt attributes
 * 5. Structured data
 * 6. URL structure
 * 7. Content optimization
 */

class SEOAnalyzer {
    constructor() {
        this.issues = [];
        this.recommendations = [];
        this.score = 0;
        this.maxScore = 100;
        
        this.initAnalysis();
    }
    
    /**
     * Initialize SEO analysis
     */
    initAnalysis() {
        window.addEventListener('load', () => {
            this.analyzePage();
            
            // Add a small delay to ensure all metrics are collected
            setTimeout(() => {
                this.displayResults();
            }, 1500);
        });
    }
    
    /**
     * Analyze the current page for SEO issues
     */
    analyzePage() {
        this.analyzeTitle();
        this.analyzeMetaDescription();
        this.analyzeHeadings();
        this.analyzeImages();
        this.analyzeStructuredData();
        this.analyzeContent();
        this.analyzeURL();
        
        // Calculate final score based on issues
        this.calculateScore();
    }
    
    /**
     * Analyze title tag
     */
    analyzeTitle() {
        const title = document.title;
        
        if (!title) {
            this.addIssue('Title tag is missing');
            return;
        }
        
        if (title.length < 30) {
            this.addIssue('Title tag is too short (less than 30 characters)');
        }
        
        if (title.length > 60) {
            this.addIssue('Title tag is too long (more than 60 characters)');
        }
        
        if (title === 'News Portal Website') {
            this.addIssue('Title tag is generic and not optimized for SEO');
        }
    }
    
    /**
     * Analyze meta description
     */
    analyzeMetaDescription() {
        const metaDescription = document.querySelector('meta[name="description"]');
        
        if (!metaDescription) {
            this.addIssue('Meta description is missing');
            return;
        }
        
        const content = metaDescription.getAttribute('content');
        
        if (!content) {
            this.addIssue('Meta description is empty');
            return;
        }
        
        if (content.length < 120) {
            this.addIssue('Meta description is too short (less than 120 characters)');
        }
        
        if (content.length > 160) {
            this.addIssue('Meta description is too long (more than 160 characters)');
        }
    }
    
    /**
     * Analyze heading structure
     */
    analyzeHeadings() {
        const h1 = document.querySelectorAll('h1');
        
        if (h1.length === 0) {
            this.addIssue('No H1 heading found on the page');
        }
        
        if (h1.length > 1) {
            this.addIssue('Multiple H1 headings found on the page');
        }
        
        const headings = document.querySelectorAll('h1, h2, h3, h4, h5, h6');
        let previousLevel = 0;
        
        for (let i = 0; i < headings.length; i++) {
            const heading = headings[i];
            const level = parseInt(heading.tagName.substring(1));
            
            if (level - previousLevel > 1 && previousLevel !== 0) {
                this.addIssue(`Heading structure skip found: H${previousLevel} to H${level}`);
            }
            
            previousLevel = level;
        }
    }
    
    /**
     * Analyze images
     */
    analyzeImages() {
        const images = document.querySelectorAll('img');
        let missingAlt = 0;
        let emptyAlt = 0;
        
        images.forEach(img => {
            if (!img.hasAttribute('alt')) {
                missingAlt++;
            } else if (img.getAttribute('alt') === '') {
                emptyAlt++;
            }
            
            // Check if image has width and height attributes
            if (!img.hasAttribute('width') || !img.hasAttribute('height')) {
                this.addIssue('Image missing width or height attributes which can cause layout shifts (CLS)');
            }
        });
        
        if (missingAlt > 0) {
            this.addIssue(`${missingAlt} images missing alt attributes`);
        }
        
        if (emptyAlt > 0) {
            this.addIssue(`${emptyAlt} images have empty alt attributes`);
        }
    }
    
    /**
     * Analyze structured data
     */
    analyzeStructuredData() {
        const structuredData = document.querySelectorAll('script[type="application/ld+json"]');
        
        if (structuredData.length === 0) {
            this.addIssue('No JSON-LD structured data found on the page');
            this.addRecommendation('Add schema.org structured data using JSON-LD for better rich snippets in search results');
        }
    }
    
    /**
     * Analyze content
     */
    analyzeContent() {
        const content = document.body.innerText;
        const words = content.split(/\s+/).filter(word => word.length > 0);
        
        if (words.length < 300) {
            this.addIssue('Content is too thin (less than 300 words)');
            this.addRecommendation('Add more relevant content to the page for better SEO results');
        }
        
        // Check if content has links
        const links = document.querySelectorAll('a');
        const internalLinks = Array.from(links).filter(link => 
            link.hostname === window.location.hostname && 
            !link.href.includes('javascript:') && 
            link.href !== window.location.href
        );
        
        if (internalLinks.length < 3) {
            this.addIssue('Page has few internal links');
            this.addRecommendation('Add more internal links to improve site structure and SEO');
        }
    }
    
    /**
     * Analyze URL structure
     */
    analyzeURL() {
        const url = window.location.pathname;
        
        // Check if URL contains .html extension
        if (url.endsWith('.html')) {
            this.addIssue('URL contains .html extension, which is not SEO-friendly');
            this.addRecommendation('Consider using clean URLs without file extensions');
        }
        
        // Check for URL length
        if (url.length > 100) {
            this.addIssue('URL is too long (over 100 characters)');
        }
        
        // Check for underscores in URL
        if (url.includes('_')) {
            this.addIssue('URL contains underscores instead of hyphens');
            this.addRecommendation('Replace underscores with hyphens in URLs for better SEO');
        }
    }
    
    /**
     * Add an issue to the list
     */
    addIssue(message) {
        this.issues.push(message);
    }
    
    /**
     * Add a recommendation
     */
    addRecommendation(message) {
        if (!this.recommendations.includes(message)) {
            this.recommendations.push(message);
        }
    }
    
    /**
     * Calculate SEO score based on issues
     */
    calculateScore() {
        // Base score of 100, subtract for each issue
        this.score = this.maxScore - (this.issues.length * 5);
        
        // Ensure score is between 0 and 100
        this.score = Math.max(0, Math.min(100, this.score));
    }
    
    /**
     * Display SEO analysis results
     */
    displayResults() {
        console.log('===== SEO Analysis Results =====');
        console.log(`SEO Score: ${this.score}/100`);
        
        if (this.issues.length > 0) {
            console.log('\nIssues found:');
            this.issues.forEach((issue, index) => {
                console.log(`${index + 1}. ${issue}`);
            });
        }
        
        if (this.recommendations.length > 0) {
            console.log('\nRecommendations:');
            this.recommendations.forEach((rec, index) => {
                console.log(`${index + 1}. ${rec}`);
            });
        }
        
        // Create visual indicator for developers
        if (localStorage.getItem('showSEOAnalysis') === 'true') {
            this.createVisualIndicator();
        }
    }
    
    /**
     * Create a visual indicator of SEO score
     */
    createVisualIndicator() {
        const indicator = document.createElement('div');
        indicator.id = 'seo-indicator';
        indicator.style.position = 'fixed';
        indicator.style.top = '20px';
        indicator.style.right = '20px';
        indicator.style.backgroundColor = this.getScoreColor(this.score);
        indicator.style.color = 'white';
        indicator.style.padding = '10px';
        indicator.style.borderRadius = '5px';
        indicator.style.zIndex = '9999';
        indicator.style.fontSize = '12px';
        indicator.style.boxShadow = '0 0 10px rgba(0, 0, 0, 0.5)';
        indicator.style.cursor = 'pointer';
        indicator.textContent = `SEO Score: ${this.score}/100`;
        
        indicator.onclick = () => {
            this.toggleAnalysisDisplay();
        };
        
        document.body.appendChild(indicator);
    }
    
    /**
     * Toggle display of full analysis
     */
    toggleAnalysisDisplay() {
        let display = document.getElementById('seo-analysis-display');
        
        if (display) {
            display.remove();
            return;
        }
        
        display = document.createElement('div');
        display.id = 'seo-analysis-display';
        display.style.position = 'fixed';
        display.style.top = '60px';
        display.style.right = '20px';
        display.style.backgroundColor = 'white';
        display.style.color = '#333';
        display.style.padding = '15px';
        display.style.borderRadius = '5px';
        display.style.zIndex = '9998';
        display.style.maxWidth = '350px';
        display.style.maxHeight = '400px';
        display.style.overflow = 'auto';
        display.style.boxShadow = '0 0 15px rgba(0, 0, 0, 0.3)';
        
        const title = document.createElement('h3');
        title.textContent = 'SEO Analysis';
        title.style.marginTop = '0';
        title.style.borderBottom = '1px solid #ddd';
        title.style.paddingBottom = '5px';
        display.appendChild(title);
        
        const score = document.createElement('p');
        score.innerHTML = `<strong>Score: ${this.score}/100</strong>`;
        score.style.color = this.getScoreColor(this.score);
        display.appendChild(score);
        
        if (this.issues.length > 0) {
            const issuesTitle = document.createElement('h4');
            issuesTitle.textContent = 'Issues:';
            display.appendChild(issuesTitle);
            
            const issuesList = document.createElement('ul');
            this.issues.forEach(issue => {
                const item = document.createElement('li');
                item.textContent = issue;
                item.style.marginBottom = '5px';
                issuesList.appendChild(item);
            });
            display.appendChild(issuesList);
        }
        
        if (this.recommendations.length > 0) {
            const recsTitle = document.createElement('h4');
            recsTitle.textContent = 'Recommendations:';
            display.appendChild(recsTitle);
            
            const recsList = document.createElement('ul');
            this.recommendations.forEach(rec => {
                const item = document.createElement('li');
                item.textContent = rec;
                item.style.marginBottom = '5px';
                recsList.appendChild(item);
            });
            display.appendChild(recsList);
        }
        
        const closeButton = document.createElement('button');
        closeButton.textContent = 'Close';
        closeButton.style.marginTop = '10px';
        closeButton.style.padding = '5px 10px';
        closeButton.style.backgroundColor = '#f0f0f0';
        closeButton.style.border = '1px solid #ddd';
        closeButton.style.borderRadius = '3px';
        closeButton.onclick = () => {
            display.remove();
        };
        display.appendChild(closeButton);
        
        document.body.appendChild(display);
    }
    
    /**
     * Get color for score display
     */
    getScoreColor(score) {
        if (score >= 80) return '#4caf50'; // Good
        if (score >= 50) return '#ff9800'; // Needs Improvement
        return '#f44336'; // Poor
    }
}

// Initialize SEO analyzer
const seoAnalyzer = new SEOAnalyzer();

// Add a method to enable/disable SEO analysis visual indicator
window.toggleSEOAnalysis = function() {
    const current = localStorage.getItem('showSEOAnalysis') === 'true';
    localStorage.setItem('showSEOAnalysis', (!current).toString());
    alert(`SEO Analysis indicator ${!current ? 'enabled' : 'disabled'}. Reload the page to see the change.`);
};
