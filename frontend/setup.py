import io

from setuptools import find_packages
from setuptools import setup

with io.open("README.md", "rt", encoding="utf8") as f:
    readme = f.read()

setup(
    name="aarontest",
    version="1.0.0",
    url="https://github.com/aarons-talks/2020-10-20-All-Things-Open",
    license="MIT",
    maintainer="Aaron Schlesinger",
    maintainer_email="aaron@ecomaz.net",
    description="Example app to show Python and Go working together",
    long_description=readme,
    packages=find_packages(),
    include_package_data=True,
    zip_safe=False,
    install_requires=["flask"],
    extras_require={"test": ["pytest", "coverage"]},
)
