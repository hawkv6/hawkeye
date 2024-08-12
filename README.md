### Analyze Command Branch

This branch introduces the `analyze` command, which is designed to generate and analyze data using various normalization methods. The resulting plots are saved in the `png` folder. This command was primarily used to evaluate and determine the most effective normalization technique.

#### Normalization Methods

The `analyze` command retrieves all link data from JAGW and normalizes them using the following methods:

- **Robust Normalizer**
- **Z-Score Normalizer**
- **MinMax Normalizer**
- **IQR MinMax Normalizer**

After testing, the **IQR MinMax Normalizer** was found to provide the best-fitting solution. Consequently, it has been implemented and tested within the `generic-processor`.

#### How to Execute

To run the `analyze` command, use the following syntax:

```bash
hawkeye analyze -j <jagw-address> -r <request-port>
