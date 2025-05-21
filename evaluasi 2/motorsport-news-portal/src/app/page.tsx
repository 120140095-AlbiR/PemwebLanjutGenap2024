export default function HomePage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-4">Motorsport News Portal</h1>
      <p className="mb-4">Welcome to our motorsport news portal.</p>
      <p>This is a simple test page to ensure the app is working correctly.</p>
    </div>
  );
}

  return (
    <div>
      {/* Hero section */}
      <section className="bg-gradient-to-r from-red-700 to-black py-12 mb-8">
        <div className="container mx-auto px-4">
          <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">
            Motorsport News Portal
          </h1>
          <p className="text-xl text-gray-200">
            Your one-stop destination for the latest news from the world of motorsports
          </p>
        </div>
      </section>
      
      {/* Filter tabs */}
      <div className="mb-8">
        <div className="flex flex-wrap gap-2">
          <Link href="/" className="bg-red-600 text-white px-4 py-2 rounded-full text-sm font-medium">
            All News
          </Link>
          {Object.values(MotorsportCategory).map((cat) => (
            <Link 
              key={cat}
              href={`/category/${encodeURIComponent(cat)}`}
              className="bg-gray-200 hover:bg-gray-300 px-4 py-2 rounded-full text-sm font-medium transition-colors"
            >
              {cat}
            </Link>
          ))}
        </div>
      </div>

      {/* News grid */}
      {loading ? (
        <div className="text-center py-20">
          <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-current border-r-transparent"></div>
          <p className="mt-4">Loading news...</p>
        </div>
      ) : error ? (
        <div className="text-center py-20 text-red-600">{error}</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {allNews.length > 0 ? (
            allNews.map((news) => (
              <NewsCard key={news.id} news={news} />
            ))
          ) : (
            <div className="col-span-3 text-center py-12">
              <p className="text-gray-600">No news articles found.</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
