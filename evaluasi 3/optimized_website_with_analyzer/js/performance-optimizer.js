/**
 * Website Performance Optimization
 * 
 * This script implements various performance optimizations:
 * 1. Lazy loading for images
 * 2. Resource hints (preload, prefetch)
 * 3. Font display optimization
 * 4. Deferred loading of non-critical resources
 * 5. Runtime performance improvements
 */

class PerformanceOptimizer {
    constructor() {
        // Initialize optimization techniques
        this.initOptimizations();
    }

    /**
     * Initialize all performance optimizations
     */
    initOptimizations() {
        window.addEventListener('load', () => {
            this.implementLazyLoading();
            this.optimizeScrollEvents();
            this.deferNonCriticalCSS();
            this.addHeightToImages();
        });
    }

    /**
     * Implement lazy loading for images that don't already have it
     */
    implementLazyLoading() {
        // Target all images that don't have loading attribute
        const images = document.querySelectorAll('img:not([loading])');
        
        images.forEach(img => {
            if (!img.hasAttribute('loading')) {
                img.setAttribute('loading', 'lazy');
                console.log('Added lazy loading to image:', img.src);
            }
        });
        
        // For browsers that don't support native lazy loading
        if ('loading' in HTMLImageElement.prototype === false) {
            this.implementIntersectionObserverLazyLoading();
        }
    }

    /**
     * Implement lazy loading using IntersectionObserver for browsers without native support
     */
    implementIntersectionObserverLazyLoading() {
        const lazyImages = document.querySelectorAll('img[loading="lazy"]');
        
        const imageObserver = new IntersectionObserver((entries, observer) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    const img = entry.target;
                    const src = img.getAttribute('data-src') || img.getAttribute('src');
                    
                    if (src) {
                        img.src = src;
                        img.removeAttribute('data-src');
                    }
                    
                    observer.unobserve(img);
                }
            });
        });
        
        lazyImages.forEach(img => {
            imageObserver.observe(img);
        });
    }

    /**
     * Optimize scroll events to prevent performance issues
     */
    optimizeScrollEvents() {
        // Use passive event listeners for scroll events
        const supportsPassive = this.checkPassiveSupport();
        
        if (supportsPassive) {
            const scrollOptions = { passive: true };
            window.addEventListener('scroll', this.throttle(this.onScroll, 100), scrollOptions);
            window.addEventListener('touchmove', this.throttle(this.onScroll, 100), scrollOptions);
        } else {
            window.addEventListener('scroll', this.throttle(this.onScroll, 100));
            window.addEventListener('touchmove', this.throttle(this.onScroll, 100));
        }
    }

    /**
     * Check if browser supports passive event listeners
     */
    checkPassiveSupport() {
        let supportsPassive = false;
        try {
            const options = {
                get passive() {
                    supportsPassive = true;
                    return true;
                }
            };
            window.addEventListener('test', null, options);
            window.removeEventListener('test', null, options);
        } catch (e) {
            console.log('Passive event listeners not supported');
        }
        return supportsPassive;
    }

    /**
     * Throttle function to limit the rate at which a function can fire
     */
    throttle(func, limit) {
        let inThrottle;
        return function() {
            const args = arguments;
            const context = this;
            if (!inThrottle) {
                func.apply(context, args);
                inThrottle = true;
                setTimeout(() => inThrottle = false, limit);
            }
        };
    }

    /**
     * Scroll event handler - can be used to implement additional optimizations
     */
    onScroll() {
        // This is intentionally left minimal for performance
        // Add any scroll-related logic here if needed
    }

    /**
     * Defer loading of non-critical CSS
     */
    deferNonCriticalCSS() {
        // Get all CSS links
        const cssLinks = document.querySelectorAll('link[rel="stylesheet"]');

        // Mark non-critical CSS for deferred loading
        const nonCriticalCSS = [
            'owl.carousel.min.css',
            'animate.min.css',
            'sweetalert2.min.css',
            'magnific-popup.css'
        ];

        cssLinks.forEach(link => {
            const href = link.getAttribute('href');
            if (href) {
                for (const css of nonCriticalCSS) {
                    if (href.includes(css)) {
                        // Convert to preload
                        link.setAttribute('rel', 'preload');
                        link.setAttribute('as', 'style');
                        link.setAttribute('onload', "this.onload=null;this.rel='stylesheet'");
                        console.log('Deferred loading for CSS:', href);
                        break;
                    }
                }
            }
        });
    }

    /**
     * Add height and width attributes to images that don't have them
     * This helps prevent layout shifts (CLS)
     */
    addHeightToImages() {
        const images = document.querySelectorAll('img:not([height])');
        
        images.forEach(img => {
            if (img.complete) {
                if (!img.hasAttribute('height') && img.naturalHeight) {
                    img.setAttribute('height', img.naturalHeight);
                    console.log('Added height attribute to image:', img.src);
                }
                
                if (!img.hasAttribute('width') && img.naturalWidth) {
                    img.setAttribute('width', img.naturalWidth);
                    console.log('Added width attribute to image:', img.src);
                }
            } else {
                img.onload = () => {
                    if (!img.hasAttribute('height') && img.naturalHeight) {
                        img.setAttribute('height', img.naturalHeight);
                    }
                    
                    if (!img.hasAttribute('width') && img.naturalWidth) {
                        img.setAttribute('width', img.naturalWidth);
                    }
                };
            }
        });
    }
}

// Initialize performance optimizer when the script loads
const performanceOptimizer = new PerformanceOptimizer();
