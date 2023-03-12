import numpy as np
from sklearn.linear_model import LinearRegression
X = np.array([1.47, 1.50, 1.52, 1.55, 1.57, 1.60, 1.63, 1.65, 1.68, 1.70, 1.73, 1.75, 1.78, 1.80, 1.83]).reshape(-1, 1)
y = np.array([52.21, 53.12, 54.48, 55.84, 57.20, 58.57, 59.93, 61.29, 63.11, 64.47, 66.28, 68.10, 69.92, 72.19, 74.46]).reshape(-1, 1)
reg = LinearRegression().fit(X, y)
print(reg.coef_, reg.intercept_, reg.score(X, y))
print(reg.predict(np.array([[1.66], [2]])))