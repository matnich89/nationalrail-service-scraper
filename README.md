The application will start multiple workers that will begin collecting train service information from the national rail api. The work is 'chunked' between workers to allow quick processing 
of the around 2.5k train stations that are in the United Kingdom. 

This is a work in progress, and was created as part of a wider project, however you are free to use this / take ideas. 

This service uses a slight modified version of https://github.com/martinsirbe/go-national-rail-client  (the copy in use is located at https://github.com/matnich89/national-rail-client)
You will need a national rail api key to run this application.

## Configuration

You can adjust the following parameters in the `main.go` file:

- `numWorkers`: Number of concurrent workers
- `maxDelay`: Maximum delay for staggered worker start


## Future Improvements

- Implement Queue interaction so service IDS can be processed / persisted by other services
- Add more comprehensive error handling and logging
- Add unit tests for critical components
- Containerize the application for easier deployment
- Clean up code structure

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License


This project is licensed under the [MIT License](LICENSE.txt) - see the [LICENSE](LICENSE.txt) file for details.